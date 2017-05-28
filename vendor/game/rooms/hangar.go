package rooms

import (
	"fmt"
	"game/db"
	"game/ent"
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sg3des/fizzgui"
)

type hangar struct {
	cMenu *fizzgui.Container
	users *fizzgui.Widget

	infoTable *informationTable

	dad       *fizzgui.DADGroup
	cShip     *fizzgui.Container
	shipSlots []*fizzgui.DADSlot
	cItems    *fizzgui.Container

	slotNormal fizzgui.Style
	slotHover  fizzgui.Style
	slotActive fizzgui.Style

	ship *ent.Ship
}

//Hangar room
func Hangar() {
	h := new(hangar)

	h.cMenu = fizzgui.NewContainer("menu", "0", "0", "100%", "10%")
	h.cMenu.AutoAdjustHeight = true

	h.users = h.cMenu.NewText("Users online: 0")
	h.users.Layout.SetWidth("80%")
	h.cMenu.NewButton("Battle", h.startBattle).Layout.SetWidth("20%")

	h.infoTable = InitInformationTable()

	//inventory
	h.dad = fizzgui.NewDragAndDropGroup("dad-group")

	h.cShip = fizzgui.NewContainer("ship", "0%", "10%", "50%", "60%")
	h.cShip.Layout.HAlign = fizzgui.HAlignRight
	h.cShip.Style.BackgroundColor = mgl32.Vec4{1, 1, 1, 1}

	h.slotNormal = fizzgui.NewStyle(fizzgui.TextColor, mgl32.Vec4{0.1, 0.2, 0.3, 0.5}, mgl32.Vec4{0.3, 0.5, 0.7, 0.5}, 1)
	h.slotHover = fizzgui.NewStyle(fizzgui.TextColor, mgl32.Vec4{0.2, 0.3, 0.4, 0.5}, mgl32.Vec4{0.3, 0.5, 0.7, 0.5}, 1)
	h.slotActive = fizzgui.NewStyle(fizzgui.TextColor, mgl32.Vec4{0.25, 0.35, 0.45, 0.6}, mgl32.Vec4{0.3, 0.5, 0.7, 0.5}, 1)

	h.cItems = fizzgui.NewContainer("items", "0%", "70%", "50%", "30%")
	h.cItems.Layout.HAlign = fizzgui.HAlignRight

	for col := 0; col < 24; col++ {
		slot := h.cItems.NewSlot(h.dad, "slot0", "", "", "12.4%", "12.4%", h.takeOff)
		slot.SetStyles(h.slotNormal, fizzgui.Style{}, h.slotActive, nil)
	}

	h.fill()
}

func (h *hangar) fill() {
	h.ship = db.GetShip("red-171014")

	h.cShip.Style.Texture = h.ship.LoadImgUI()
	for _, s := range h.ship.Slots {
		slot := h.cShip.NewSlot(h.dad, string(s.ET), s.X, s.Y, s.W, s.H, h.putOn)
		slot.SetStyles(h.slotNormal, h.slotHover, h.slotActive, nil)
		h.shipSlots = append(h.shipSlots, slot)
	}

	weapons := []string{"gun-104", "laser-89"}
	for i, model := range weapons {
		w := db.GetWeapon(model)
		item := h.dad.NewItem(w.ET.Str(), w.ImgPath(), w)
		h.dad.Slots[i].PlaceItem(item)
	}

	items := []string{"engine-22"}
	for i, model := range items {
		e := db.GetEquipment(model)
		item := h.dad.NewItem(e.ET.Str(), e.ImgPath(), e)
		h.dad.Slots[i+len(weapons)].PlaceItem(item)
	}

	h.recalculate()
}

func (h *hangar) putOn(item *fizzgui.DADItem, slot *fizzgui.DADSlot, value interface{}) bool {
	if item.ID != slot.ID {
		//ID contains type of equipment, item possible to place in equal type, ex: weapon to weapon slot
		return false
	}

	slot.PlaceItem(item)
	h.recalculate()

	return true
}

func (h *hangar) takeOff(item *fizzgui.DADItem, slot *fizzgui.DADSlot, value interface{}) bool {

	slot.PlaceItem(item)
	h.recalculate()

	return true
}

func (h *hangar) recalculate() {
	h.ship.Equipment = nil
	h.ship.Weapons = nil

	h.ship.Spec = h.ship.ShipSpec

	for _, slot := range h.dad.Slots {
		if slot.Item == nil {
			continue
		}

		switch slot.ID {
		case string(ent.ETweapon):
			h.ship.Weapons = append(h.ship.Weapons, slot.Item.Value.(*ent.Weapon))
		case string(ent.ETengine):
			h.ship.Equipment = append(h.ship.Equipment, slot.Item.Value.(*ent.Equipment))
		}
	}

	for _, e := range h.ship.Equipment {
		h.ship.Spec = h.ship.Spec.Summ(e.Spec)
	}

	for _, w := range h.ship.Weapons {
		h.ship.Spec.Weight += w.Weight
	}

	h.infoTable.Class.update(h.ship.Class, h.ship.Model)
	h.infoTable.Weight.update(h.ship.Spec.Weight)
	h.infoTable.Durability.update(h.ship.Spec.Durability)
	h.infoTable.Shield.update(h.ship.Spec.Shield)
	h.infoTable.Speed.update(h.ship.Spec.Speed)
	h.infoTable.Contollability.update(h.ship.Spec.Controllability)
	h.infoTable.Energy.update(h.ship.Spec.Energy, h.ship.Spec.EnergyAcc)
	h.infoTable.Metal.update(h.ship.Spec.Metal, h.ship.Spec.MetalAcc)

	for i, w := range h.ship.Weapons {
		h.infoTable.newInfoWeapon(i, w)
	}

	for i := len(h.ship.Weapons); i < len(h.infoTable.Weapons); i++ {
		h.infoTable.Weapons[i].SetHidden(true)
	}

	log.Println(db.Ship.Weapons)
}

func (h *hangar) startBattle(wgt *fizzgui.Widget) {

}

type informationTable struct {
	C *fizzgui.Container

	Class          *field
	Weight         *field
	Durability     *field
	Shield         *field
	Speed          *field
	Contollability *field
	Energy         *field
	Metal          *field

	Weapons []*infoWeapon
}

type field struct {
	widget *fizzgui.Widget
	format string
	vals   []interface{}
}

func (info *informationTable) newField(format string) *field {
	return &field{
		widget: info.newText(""),
		format: format,
	}
}

func (f *field) update(vals ...interface{}) {
	f.vals = vals
	f.widget.Text = fmt.Sprintf(f.format, vals...)
}

func InitInformationTable() *informationTable {
	info := new(informationTable)
	info.C = fizzgui.NewContainer("info", "50%", "10%", "30%", "90%")
	info.C.Layout.HAlign = fizzgui.HAlignRight

	info.Class = info.newField("Class: %s [%s]")
	info.Weight = info.newField("Weight:         %.1f")
	info.Durability = info.newField("Durability:     %.0f")
	info.Shield = info.newField("Shield:         %.0f")
	info.Speed = info.newField("Speed:          %.1f")
	info.Contollability = info.newField("Contollability: %.1f")
	info.Energy = info.newField("Energy:         %.0f [+%.1f]")
	info.Metal = info.newField("Metal:          %.0f [+%.1f]")

	return info
}

type infoWeapon struct {
	Name  *fizzgui.Widget
	DPS   *fizzgui.Widget
	Range *fizzgui.Widget
	Ammo  *fizzgui.Widget
}

func (info *informationTable) newInfoWeapon(i int, w *ent.Weapon) {
	dps := fmt.Sprintf("  DPS: %.1f", w.Damage/float32(w.AttackRate.Seconds()))
	rang := fmt.Sprintf("  Range: %.1f", w.AttackRange)
	ammo := fmt.Sprintf("  Ammo: %d", w.Ammunition)

	if i < len(info.Weapons) {
		iw := info.Weapons[i]
		iw.SetHidden(false)
		iw.Name.Text = w.Model
		iw.DPS.Text = dps
		iw.Range.Text = rang
		iw.Ammo.Text = ammo
	} else {
		iw := &infoWeapon{
			Name:  info.newText(w.Model),
			DPS:   info.newText(dps),
			Range: info.newText(rang),
			Ammo:  info.newText(ammo),
		}
		info.Weapons = append(info.Weapons, iw)
	}

	//name string, damage, rate, rang float32, ammo int

	// info.Weapons = append(info.Weapons, w)
}

func (iw *infoWeapon) SetHidden(b bool) {
	iw.Name.Hidden = b
	iw.DPS.Hidden = b
	iw.Range.Hidden = b
	iw.Ammo.Hidden = b
}

func (info *informationTable) newText(str string) *fizzgui.Widget {
	w := info.C.NewText(str)
	w.Layout.SetWidth("100%")
	w.Layout.Margin = fizzgui.Offset{2, 2, 2, 2}
	w.Layout.Padding = fizzgui.Offset{2, 2, 2, 2}
	w.Font = fizzgui.GetFont("Mono")
	return w
}
