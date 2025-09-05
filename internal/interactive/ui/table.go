package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/web-seven/overlock/internal/interactive/ui/theme"
)

func (m *Model) updateTable() {
	tableWidth := m.windowWidth - theme.WindowPadding
	if tableWidth <= 0 {
		tableWidth = theme.DefaultWidth
	}

	if tableWidth < theme.MinTerminalWidth {
		tableWidth = theme.MinTerminalWidth
	}

	availableWidth := tableWidth - theme.TableBorderSpace

	nameWidth := max(theme.NameMinWidth, min(availableWidth*theme.NameWidthPercent/100, theme.NameMaxWidth))
	packageWidth := max(theme.PackageMinWidth, min(availableWidth*theme.PackageWidthPercent/100, theme.PackageMaxWidth))
	versionWidth := max(theme.VersionMinWidth, min(availableWidth*theme.VersionWidthPercent/100, theme.VersionMaxWidth))
	statusWidth := max(theme.StatusMinWidth, min(availableWidth*theme.StatusWidthPercent/100, theme.StatusMaxWidth))
	dateWidth := max(theme.DateMinWidth, min(availableWidth*theme.DateWidthPercent/100, theme.DateMaxWidth))

	usedWidth := nameWidth + packageWidth + versionWidth + statusWidth + dateWidth
	descWidth := max(theme.DescMinWidth, availableWidth-usedWidth)

	columns := []table.Column{
		{Title: "Name", Width: nameWidth},
		{Title: "Package", Width: packageWidth},
		{Title: "Version", Width: versionWidth},
		{Title: "Status", Width: statusWidth},
		{Title: "Install Date", Width: dateWidth},
		{Title: "Description", Width: descWidth},
	}

	var rows []table.Row

	// Get current tab state
	state := m.tabStates[m.activeTab]
	items := state.items

	switch m.activeTab {
	case ConfigurationsTab:
		rows = make([]table.Row, len(items))
		for i, config := range items {
			name := truncateString(config.Name, nameWidth-theme.CellPadding)
			pkg := truncateString(config.Package, packageWidth-theme.CellPadding)
			version := truncateString(config.Version, versionWidth-theme.CellPadding)
			status := truncateString(config.Status, statusWidth-theme.CellPadding)
			date := truncateString(config.InstallDate, dateWidth-theme.CellPadding)
			desc := truncateString(config.Description, descWidth-theme.CellPadding)
			rows[i] = table.Row{name, pkg, version, status, date, desc}
		}
	case ProvidersTab:
		rows = make([]table.Row, len(items))
		for i, provider := range items {
			name := truncateString(provider.Name, nameWidth-theme.CellPadding)
			pkg := truncateString(provider.Package, packageWidth-theme.CellPadding)
			version := truncateString(provider.Version, versionWidth-theme.CellPadding)
			status := truncateString(provider.Status, statusWidth-theme.CellPadding)
			date := truncateString(provider.InstallDate, dateWidth-theme.CellPadding)
			desc := truncateString(provider.Description, descWidth-theme.CellPadding)
			rows[i] = table.Row{name, pkg, version, status, date, desc}
		}
	case FunctionsTab:
		rows = make([]table.Row, len(items))
		for i, function := range items {
			name := truncateString(function.Name, nameWidth-theme.CellPadding)
			pkg := truncateString(function.Package, packageWidth-theme.CellPadding)
			version := truncateString(function.Version, versionWidth-theme.CellPadding)
			status := truncateString(function.Status, statusWidth-theme.CellPadding)
			date := truncateString(function.InstallDate, dateWidth-theme.CellPadding)
			desc := truncateString(function.Description, descWidth-theme.CellPadding)
			rows[i] = table.Row{name, pkg, version, status, date, desc}
		}
	default:
		rows = []table.Row{}
	}

	m.table.SetColumns(columns)
	m.table.SetRows(rows)
}

func truncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}

	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\v", " ")
	s = strings.ReplaceAll(s, "\f", " ")

	s = strings.TrimSpace(strings.Join(strings.Fields(s), " "))

	if s == "" {
		return ""
	}

	if maxLen <= theme.MinEllipsisWidth {
		if len(s) <= maxLen {
			return s
		}
		return s[:maxLen]
	}

	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-theme.EllipsisLength] + "..."
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
