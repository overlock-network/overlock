package resources

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/web-seven/overlock/internal/engine"
	"github.com/web-seven/overlock/internal/kube"
	"go.uber.org/zap"

	crossv1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	v1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type XResource struct {
	Resource string
	unstructured.Unstructured
}

var apiFields = []string{"apiVersion", "kind"}
var metadataFields = []string{"metadata"}

func (xr *XResource) GetSchemaFormFromXRDefinition(ctx context.Context, xrd crossv1.CompositeResourceDefinition, client *dynamic.DynamicClient, logger *zap.SugaredLogger) *huh.Form {

	xrdInstance, err := client.Resource(schema.GroupVersionResource{
		Group:    xrd.GroupVersionKind().Group,
		Version:  xrd.GroupVersionKind().Version,
		Resource: xrd.GroupVersionKind().Kind,
	}).Get(ctx, xrd.Name, metav1.GetOptions{})

	if err != nil {
		logger.Error(err)
		return nil
	}

	runtime.DefaultUnstructuredConverter.FromUnstructured(xrdInstance.UnstructuredContent(), &xrd)

	formGroups := []*huh.Group{}

	selectedVersion := v1.CompositeResourceDefinitionVersion{}
	if len(xrd.Spec.Versions) == 1 {
		selectedVersion = xrd.Spec.Versions[0]
	} else {
		selectedVersionIndex := 0
		versionOptions := []huh.Option[int]{}
		for index, version := range xrd.Spec.Versions {
			versionOptions = append(versionOptions, huh.NewOption(version.Name, index))
		}
		vesionSelectForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[int]().
					Title("Version").
					Options(versionOptions...).
					Value(&selectedVersionIndex),
			),
		)
		vesionSelectForm.Run()
		selectedVersion = xrd.Spec.Versions[selectedVersionIndex]
	}

	versionSchema, _ := parseSchema(selectedVersion.Schema, logger)

	logger.Info("Type: \t\t" + xrd.Name)
	logger.Info("Description: \t" + versionSchema.Description)

	versionGroups := xr.getFormGroupsByProps(versionSchema, "")
	formGroups = append(formGroups, versionGroups...)

	xr.Unstructured.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   xrd.Spec.Group,
		Version: selectedVersion.Name,
	})

	xr.Unstructured.SetKind(xrd.Spec.Names.Kind)
	xr.Resource = xrd.Spec.Names.Plural

	formGroups = append(formGroups,
		huh.NewGroup(
			huh.NewConfirm().
				Key("confirm").
				Title("Would you like to create resource?"),
		),
	)

	schemaForm := huh.NewForm(formGroups...)

	return schemaForm
}

func (xr *XResource) getFormGroupsByProps(schema *extv1.JSONSchemaProps, parent string) []*huh.Group {
	formGroups := []*huh.Group{}
	formFields := []huh.Field{}

	if xr.Unstructured.Object == nil {
		xr.Unstructured.Object = make(map[string]interface{})
	}
	for propertyName, property := range schema.Properties {

		breadCrumbs := parent + "[" + propertyName + "]"

		shortDescription := strings.SplitN(property.Description, ".", 2)[0]
		description := shortDescription + breadCrumbs
		isRequired := isStringInArray(schema.Required, propertyName)

		if (property.Type == "object" || property.Type == "array") && (!isStringInArray(metadataFields, propertyName) || len(property.Properties) > 0) {
			schemaXr := XResource{}
			(xr.Unstructured.Object)[propertyName] = &schemaXr.Unstructured.Object

			if property.Type == "array" {
				if property.Items.Schema.Type == "string" {

					propertyValue := []string{}

					if property.Items.Schema.Description != "" {
						description = strings.SplitN(property.Items.Schema.Description, ".", 2)[0]
					}
					enums := []string{}
					for _, optionValue := range property.Items.Schema.Enum {
						timmedValue := strings.Trim(string(optionValue.Raw), "\"")
						enums = append(enums, timmedValue)
					}

					formFields = append(formFields,
						huh.NewText().
							Title(description).
							Lines(3).
							Validate(func(s string) error {

								if s != "" {
									if len(enums) > 0 {
										propertyValues := strings.Split(s, "\n")
										for _, optionValue := range propertyValues {
											if !isStringInArray(enums, optionValue) {
												return errors.New("supported values: " + strings.Join(enums, ", "))
											}
										}
									}
									propertyValue = strings.Split(s, "\n")

								}
								return nil
							}),
					)
					(xr.Unstructured.Object)[propertyName] = &propertyValue
				} else if property.Items.Schema.Type == "object" {
					propertyGroups := schemaXr.getFormGroupsByProps(property.Items.Schema, breadCrumbs)
					formGroups = append(formGroups, propertyGroups...)
					(xr.Unstructured.Object)[propertyName] = &[]map[string]interface{}{schemaXr.Unstructured.Object}

				}

			} else {
				propertyGroups := schemaXr.getFormGroupsByProps(&property, breadCrumbs)
				formGroups = append(formGroups, propertyGroups...)
				(xr.Unstructured.Object)[propertyName] = &schemaXr.Unstructured.Object
			}

		} else if property.Type == "string" && !isStringInArray(apiFields, propertyName) {
			propertyValue := ""
			(xr.Unstructured.Object)[propertyName] = &propertyValue

			if len(property.Enum) > 0 {
				if property.Default != nil {
					propertyValue = strings.Trim(string(property.Default.Raw), "\"")
				}
				options := []huh.Option[string]{}
				for _, optionValue := range property.Enum {
					timmedValue := strings.Trim(string(optionValue.Raw), "\"")
					options = append(options, huh.NewOption(timmedValue, timmedValue))
				}
				formFields = append(formFields, huh.NewSelect[string]().
					Options(options...).
					Title(description).
					Value(&propertyValue),
				)
			} else {
				formFields = append(formFields, huh.NewInput().
					Validate(func(s string) error {
						if isRequired && s == "" {
							return errors.New(propertyName + " is required")
						} else {
							return nil
						}
					}).
					Title(description).
					Value(&propertyValue),
				)
			}

		} else if property.Type == "number" && !isStringInArray(apiFields, propertyName) {
			propertyValue := json.Number("")
			(xr.Unstructured.Object)[propertyName] = &propertyValue
			formFields = append(formFields, huh.NewInput().
				Validate(func(s string) error {

					if s != "" && !regexp.MustCompile(`\d`).MatchString(s) {
						return errors.New(propertyName + " shall be numeric")
					}

					if isRequired && s == "" {
						return errors.New(propertyName + " is required")
					}

					propertyValue = json.Number(s)
					return nil

				}).
				Title(description),
			)
		} else if property.Type == "boolean" && !isStringInArray(apiFields, propertyName) {
			propertyValue := false
			(xr.Unstructured.Object)[propertyName] = &propertyValue
			formFields = append(formFields, huh.NewConfirm().
				Title(description).
				Value(&propertyValue),
			)
		} else if property.Type == "object" && isStringInArray(metadataFields, propertyName) {
			propertyValue := metav1.ObjectMeta{
				Name: "",
				Labels: map[string]string{
					"app.kubernetes.io/managed-by": "overlock",
					"creation-date":                time.Now().String(),
					"update-date":                  time.Now().String(),
				},
			}
			(xr.Unstructured.Object)[propertyName] = &propertyValue

			formFields = append(formFields, huh.NewInput().
				Validate(func(s string) error {
					if s == "" {
						return errors.New("name is required")
					} else {
						return nil
					}
				}).
				Title("Name of resource").
				Value(&propertyValue.Name),
			)
		}

	}
	if len(formFields) > 0 {
		group := huh.NewGroup(formFields...).Description(schema.Description)
		currentGroups := []*huh.Group{}
		currentGroups = append(currentGroups, group)
		formGroups = append(currentGroups, formGroups...)
	}

	return formGroups
}

func parseSchema(v *v1.CompositeResourceValidation, logger *zap.SugaredLogger) (*extv1.JSONSchemaProps, error) {
	if v == nil {
		return nil, nil
	}

	s := &extv1.JSONSchemaProps{}
	if err := json.Unmarshal(v.OpenAPIV3Schema.Raw, s); err != nil {
		logger.Error(err)
	}
	return s, nil
}

func isStringInArray(a []string, s string) bool {
	for _, e := range a {
		if s == e {
			return true
		}
	}
	return false
}

func ApplyResources(ctx context.Context, client *dynamic.DynamicClient, logger *zap.SugaredLogger, file string) error {
	resources, err := transformToUnstructured(file, logger)

	if err != nil {
		return err
	}
	for _, resource := range resources {
		apiAndVersion := strings.Split(resource.GetAPIVersion(), "/")

		resourceId := schema.GroupVersionResource{
			Group:    apiAndVersion[0],
			Version:  apiAndVersion[1],
			Resource: strings.ToLower(resource.GetKind()) + "s",
		}
		resource.SetLabels(engine.ManagedLabels(nil))
		logger.Infof("Applying resource: %s", resourceId.String())
		res, err := client.Resource(resourceId).Apply(ctx, resource.GetName(), &resource, metav1.ApplyOptions{FieldManager: "overlock"})

		if err != nil {
			return err
		} else {
			logger.Infof("Resource %s from %s successfully applied", res.GetName(), res.GetAPIVersion())
		}
	}
	return nil
}

func CopyComposites(ctx context.Context, logger *zap.SugaredLogger, sourceContext dynamic.Interface, destinationContext dynamic.Interface) error {

	//Get composite resources from XRDs definition and apply them
	XRDs, err := kube.GetKubeResources(kube.ResourceParams{
		Dynamic:    sourceContext,
		Ctx:        ctx,
		Group:      "apiextensions.crossplane.io",
		Version:    "v1",
		Resource:   "compositeresourcedefinitions",
		Namespace:  "",
		ListOption: metav1.ListOptions{},
	})
	if err != nil {
		return err
	}

	if len(XRDs) > 0 {
		for _, xrd := range XRDs {
			var paramsXRs v1.CompositeResourceDefinition
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(xrd.UnstructuredContent(), &paramsXRs); err != nil {
				logger.Infof("Failed to convert item %s: %v\n", xrd.GetName(), err)
				return nil
			}
			for _, version := range paramsXRs.Spec.Versions {
				XRs, err := kube.GetKubeResources(kube.ResourceParams{
					Dynamic:   sourceContext,
					Ctx:       ctx,
					Group:     paramsXRs.Spec.Group,
					Version:   version.Name,
					Resource:  paramsXRs.Spec.Names.Plural,
					Namespace: "",
					ListOption: metav1.ListOptions{
						LabelSelector: "app.kubernetes.io/managed-by=overlock",
					},
				})
				if err != nil {
					return err
				}

				for _, xr := range XRs {
					xr.SetResourceVersion("")
					xr.SetFinalizers(nil)
					resourceId := schema.GroupVersionResource{
						Group:    paramsXRs.Spec.Group,
						Version:  version.Name,
						Resource: paramsXRs.Spec.Names.Plural,
					}

					_, err = destinationContext.Resource(resourceId).Namespace("").Get(ctx, xr.GetName(), metav1.GetOptions{})
					if err != nil {
						_, err = destinationContext.Resource(resourceId).Namespace("").Create(ctx, &xr, metav1.CreateOptions{})
						if err != nil {
							logger.Warn(err)
						} else {
							logger.Infof("Resource created successfully %s", xr.GetName())
						}
					} else {
						logger.Warnf("Resource %s with type %s already exists, skipping.", xr.GetName(), resourceId.GroupResource().String())
					}
				}
			}
		}
	} else {
		logger.Warn("Composite resources not found")
	}
	return nil
}
