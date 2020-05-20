package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*Menu)(nil)
var _ fyne.Tappable = (*Menu)(nil)

// Menu is a widget for displaying a fyne.Menu.
type Menu struct {
	Base
	DismissAction func()
	Items         []fyne.CanvasObject
	activeChild   *Menu
	customSized   bool
}

// NewMenu creates a new Menu.
func NewMenu(menu *fyne.Menu) *Menu {
	items := make([]fyne.CanvasObject, len(menu.Items))
	m := &Menu{Items: items}
	for i, item := range menu.Items {
		if item.IsSeparator {
			items[i] = NewMenuItemSeparator()
		} else {
			items[i] = NewMenuItem(item, m)
		}
	}
	return m
}

// CreateRenderer returns a new renderer for the menu.
// Implements: fyne.Widget
func (m *Menu) CreateRenderer() fyne.WidgetRenderer {
	cont := &fyne.Container{
		Layout:  layout.NewVBoxLayout(),
		Objects: m.Items,
	}
	return &menuRenderer{
		NewShadowingRenderer([]fyne.CanvasObject{cont}, MenuLevel),
		cont,
		m,
	}
}

// DeactivateChild deactivates the active child menu.
func (m *Menu) DeactivateChild() {
	if m.activeChild != nil {
		m.activeChild.Hide()
		m.activeChild = nil
	}
}

// Hide hides the menu.
// Implements: fyne.Widget
func (m *Menu) Hide() {
	HideWidget(&m.Base, m)
}

// MinSize returns the minimal size of the menu.
// Implements: fyne.Widget
func (m *Menu) MinSize() fyne.Size {
	return MinSizeOf(m)
}

// Refresh triggers a redraw of the menu.
// Implements: fyne.Widget
func (m *Menu) Refresh() {
	RefreshWidget(m)
}

// Resize has no effect because menus are always displayed with their minimal size.
// Implements: fyne.Widget
func (m *Menu) Resize(size fyne.Size) {
	ResizeWidget(&m.Base, m, size)
}

// Show makes the menu visible.
// Implements: fyne.Widget
func (m *Menu) Show() {
	ShowWidget(&m.Base, m)
}

// Tapped catches taps on separators and the menu background. It doesn’t perform any action.
// Implements: fyne.Tappable
func (m *Menu) Tapped(*fyne.PointEvent) {
	// Hit a separator or padding -> do nothing.
}

// Dismiss dismisses the menu by dismissing and hiding the active child and performing the DismissAction.
func (m *Menu) Dismiss() {
	if m.activeChild != nil {
		defer m.activeChild.Dismiss()
		m.activeChild.Hide()
		m.activeChild = nil
	}
	if m.DismissAction != nil {
		m.DismissAction()
	}
}

type menuRenderer struct {
	*ShadowingRenderer
	cont *fyne.Container
	m    *Menu
}

func (r *menuRenderer) Layout(s fyne.Size) {
	minSize := r.MinSize()
	var size fyne.Size
	if r.m.customSized {
		size = minSize.Max(s)
	} else {
		size = minSize
	}
	if size != r.m.Size() {
		r.m.Resize(size)
		return
	}

	r.LayoutShadow(size, fyne.NewPos(0, 0))
	padding := r.padding()
	r.cont.Resize(size.Subtract(padding))
	r.cont.Move(fyne.NewPos(padding.Width/2, padding.Height/2))
}

func (r *menuRenderer) MinSize() fyne.Size {
	return r.cont.MinSize().Add(r.padding())
}

func (r *menuRenderer) Refresh() {
	canvas.Refresh(r.m)
}

func (r *menuRenderer) padding() fyne.Size {
	return fyne.NewSize(0, theme.Padding()*2)
}
