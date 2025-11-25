package kube

import (
	"context"
	"fmt"
	"time"

	"github.com/pterm/pterm"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	DefaultAdminServiceAccountName = "overlock-admin"
	ClusterAdminRole               = "cluster-admin"
)

// AdminServiceAccountInfo contains the created service account information
type AdminServiceAccountInfo struct {
	Name      string
	Namespace string
	Token     string
}

// CreateAdminServiceAccount creates a service account with cluster-admin privileges
func CreateAdminServiceAccount(ctx context.Context, config *rest.Config, serviceAccountName, targetNamespace string, logger *zap.SugaredLogger) (*AdminServiceAccountInfo, error) {
	if serviceAccountName == "" {
		serviceAccountName = DefaultAdminServiceAccountName
	}

	client, err := Client(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	logger.Infof("Creating admin service account '%s' in namespace '%s'", serviceAccountName, targetNamespace)

	// Create service account
	sa, err := createServiceAccount(ctx, client, serviceAccountName, targetNamespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create service account: %w", err)
	}

	// Create cluster role binding
	err = createClusterRoleBinding(ctx, client, serviceAccountName, targetNamespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster role binding: %w", err)
	}

	// Generate token
	token, err := createServiceAccountToken(ctx, client, serviceAccountName, targetNamespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create service account token: %w", err)
	}

	info := &AdminServiceAccountInfo{
		Name:      sa.Name,
		Namespace: sa.Namespace,
		Token:     token,
	}

	displayServiceAccountInfo(info, logger)

	return info, nil
}

// createServiceAccount creates a new service account
func createServiceAccount(ctx context.Context, client *kubernetes.Clientset, name, namespace string) (*corev1.ServiceAccount, error) {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "overlock",
				"app.kubernetes.io/component":  "admin-service-account",
			},
		},
	}

	createdSA, err := client.CoreV1().ServiceAccounts(namespace).Create(ctx, sa, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		return nil, err
	}

	if errors.IsAlreadyExists(err) {
		// Get existing service account
		createdSA, err = client.CoreV1().ServiceAccounts(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
	}

	return createdSA, nil
}

// createClusterRoleBinding creates a cluster role binding for cluster-admin access
func createClusterRoleBinding(ctx context.Context, client *kubernetes.Clientset, serviceAccountName, namespace string) error {
	crbName := fmt.Sprintf("%s-cluster-admin", serviceAccountName)
	
	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: crbName,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "overlock",
				"app.kubernetes.io/component":  "admin-service-account",
			},
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     ClusterAdminRole,
		},
	}

	_, err := client.RbacV1().ClusterRoleBindings().Create(ctx, crb, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

// createServiceAccountToken creates a token for the service account
func createServiceAccountToken(ctx context.Context, client *kubernetes.Clientset, serviceAccountName, namespace string) (string, error) {
	secretName := fmt.Sprintf("%s-token", serviceAccountName)
	
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Annotations: map[string]string{
				"kubernetes.io/service-account.name": serviceAccountName,
			},
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "overlock",
				"app.kubernetes.io/component":  "admin-service-account",
			},
		},
		Type: corev1.SecretTypeServiceAccountToken,
	}

	_, err := client.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		return "", err
	}

	// Wait for token to be populated
	for i := 0; i < 30; i++ {
		secret, err := client.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}

		if token, exists := secret.Data["token"]; exists && len(token) > 0 {
			return string(token), nil
		}

		time.Sleep(1 * time.Second)
	}

	return "", fmt.Errorf("timeout waiting for service account token to be populated")
}

// displayServiceAccountInfo displays the service account information to the user
func displayServiceAccountInfo(info *AdminServiceAccountInfo, logger *zap.SugaredLogger) {
	logger.Info("Admin service account created successfully!")
	
	// Security warning
	pterm.Warning.Println("This service account has cluster-admin privileges.")
}

