package client

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func displayLoop() {
	app := tview.NewApplication()

	// Header (Application Name and Version)
	header := tview.NewTextView().
		SetText("[white]🔍 Server Monitor v1.0 - by Hany Mamdouh").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	// Footer (Instructions)
	footer := tview.NewTextView().
		SetText("[white] [Press 'q' to exit] [Press 'enter' to focus details] [Press 'Esc' to focus table]").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	// Create the table for the upper panel
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false) // Enable row selection

	// Customize highlight style to change only text color
	table.SetSelectedStyle(tcell.StyleDefault.
		Foreground(tcell.ColorYellow)) // Highlight text without changing background

	// Define headers
	headers := []string{"Hostname", "Memory (Free/Total)", "CPU Load", "Monitored Processes"}

	// Create a text view for details
	detailsView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetScrollable(true).
		SetText("Select a server to view details.")

	// Track scroll position for detailsView
	scrollOffset := 0
	// Keyboard handling for detailsView
	detailsView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc: // Escape key returns focus to table
			app.SetFocus(table)
			return nil
		case tcell.KeyUp:
			scrollOffset-- // Scroll up
			detailsView.ScrollTo(scrollOffset, 0)
			return nil
		case tcell.KeyDown:
			scrollOffset++ // Scroll down
			detailsView.ScrollTo(scrollOffset, 0)
			return nil
		case tcell.KeyPgUp:
			scrollOffset -= 5 // Scroll page up
			detailsView.ScrollTo(scrollOffset, 0)
			return nil
		case tcell.KeyPgDn:
			scrollOffset += 5 // Scroll page down
			detailsView.ScrollTo(scrollOffset, 0)
			return nil
		}
		return event
	})

	// Function to update the table
	updateTable := func() {
		mu.Lock()
		defer mu.Unlock()

		table.Clear()

		// Set headers
		for i, header := range headers {
			table.SetCell(0, i, tview.NewTableCell(header).
				SetAlign(tview.AlignCenter).
				SetSelectable(false).
				SetExpansion(1))
		}

		// Populate table with data
		for i, data := range serverDataList {
			table.SetCell(i+1, 0, tview.NewTableCell(data.HostName).SetAlign(tview.AlignCenter))
			table.SetCell(i+1, 1, tview.NewTableCell(fmt.Sprintf("%d / %dMB", data.FreeMem/1024/1024, data.TotalMemory/1024/1024)).SetAlign(tview.AlignCenter))
			table.SetCell(i+1, 2, tview.NewTableCell(fmt.Sprintf("%d%%", data.CPULoad)).SetAlign(tview.AlignCenter))
			table.SetCell(i+1, 3, tview.NewTableCell(fmt.Sprintf("%d", len(data.Processes))).SetAlign(tview.AlignCenter))
		}
	}

	// Function to update details view
	updateDetails := func(row int) {
		mu.Lock()
		defer mu.Unlock()

		if row > 0 && row <= len(serverDataList) {
			data := serverDataList[row-1]
			details := fmt.Sprintf(
				"[green]Hostname:[-] %s\n"+
					"[green]Free Memory:[-] %d MB\n"+
					"[green]CPU Load:[-] %d%%\n"+
					"[green]Disk Usage:[-]\n%s\n"+
					"[green]Processes:[-]\n",
				data.HostName, data.FreeMem/1024/1024, data.CPULoad, renderDiskUsage(data.DiskUsage),
			)

			for _, proc := range data.Processes {
				details += fmt.Sprintf(
					"  [yellow]Name:[-] %s\n"+
						"  [yellow]Status:[-] %v\n"+
						"  [yellow]Memory Usage:[-] %d KB\n"+
						"  [yellow]CPU Usage:[-] %.2f%%\n"+
						"  [yellow]Process ID:[-] %d\n"+
						"  [yellow]Logs:[-] %s\n\n",
					proc.Name, proc.Status, proc.MemoryUsage/1024, proc.CPUUsage, proc.ProcessID, proc.ProcessLogs,
				)
			}

			detailsView.SetText(details)
		} else {
			detailsView.SetText("No server selected.")
		}
	}

	// Initial table update
	updateTable()

	// Table selection callback
	table.SetSelectedFunc(func(row, _ int) {
		if row > 0 {
			updateDetails(row)
			app.SetFocus(detailsView) // Automatically switch focus to detailsView when a row is selected
		}
	})

	// Default selection
	table.Select(1, 0)
	updateDetails(1)

	// Main content (table + details)
	mainContent := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true).
		AddItem(detailsView, 0, 1, false)

	// Layout (Header + Main Content + Footer)
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 2, 0, false). // Fixed height for header
		AddItem(mainContent, 0, 1, true).
		AddItem(footer, 2, 0, false) // Fixed height for footer

	// Refresh table every 2 seconds
	go func() {
		for {
			time.Sleep(2 * time.Second)
			app.QueueUpdateDraw(func() {
				updateTable()
				row, _ := table.GetSelection()
				updateDetails(row)
			})
		}
	}()

	// Capture 'q' key to exit
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			app.Stop()
			os.Exit(0)
		}
		return event
	})

	// Start application
	if err := app.SetRoot(layout, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}

func renderDiskUsage(diskUsage map[string]int) string {
	// Find the longest partition name
	maxLen := 0
	for partition := range diskUsage {
		if len(partition) > maxLen {
			maxLen = len(partition)
		}
	}

	// Generate formatted disk usage output
	var output string
	for partition, usage := range diskUsage {
		// Pad partition names with spaces to align bars
		paddedPartition := fmt.Sprintf("%-*s", maxLen, partition)

		// Determine color based on usage
		barColor := "[green]" // Default to green
		if usage >= 75 {
			barColor = "[red]" // High usage
		} else if usage >= 50 {
			barColor = "[yellow]" // Medium usage
		}

		// Generate progress bar
		progressBar := strings.Repeat("█", usage/10) + strings.Repeat("░", 10-(usage/10)) // 10-block bar

		output += fmt.Sprintf("[white]%s:[-] %s%s[-] %d%%\n", paddedPartition, barColor, progressBar, usage)
	}
	return output
}
