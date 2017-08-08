package rooms

import (
	"fmt"
	"game"
	"game/db"
	"game/equip"
	"game/weapons"
	"log"
	"path/filepath"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sg3des/fizzgui"
)

var player *game.Player

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
}

//Hangar room
func Hangar(playerName string) {
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

	player = db.GetPlayer(playerName)
	// ship = db.GetHangarShip("red-171014")

	h.fill()
}

func LoadImage(dir, imgname string) *fizzgui.Texture {
	filename := ImgPath(dir, imgname)
	tex, err := fizzgui.NewTextureImg(filename)
	if err != nil {
		log.Fatalln("texture '%s' not found", filename)
		return nil
	}
	return tex
}

func ImgPath(dir, imgname string) string {
	return filepath.Join("assets", dir, imgname+".png")
}

func (h *hangar) fill() {
	h.cShip.Style.Texture = LoadImage("ships", player.Ship.Img)
	for _, s := range player.Ship.Slots {
		slot := h.cShip.NewSlot(h.dad, s.Type.Str(), s.X, s.Y, s.W, s.H, h.putOn)
		slot.UserData = s
		slot.SetStyles(h.slotNormal, h.slotHover, h.slotActive, nil)
		h.shipSlots = append(h.shipSlots, slot)
	}

	var placed int
	for i, w := range db.GetPlayerWeapons(player.WeaponsIDs) {
		item := h.dad.NewItem(w.EquipType.Str(), ImgPath("items", w.Img), w)

		if slot, ok := h.lookupShipSlotByName(w.SlotName); ok {
			slot.PlaceItem(item)
		} else {
			h.dad.Slots[i].PlaceItem(item)
			placed++
		}
	}

	for i, e := range db.GetPlayerEquip(player.EquipIDs) {
		item := h.dad.NewItem(e.Type.Str(), ImgPath("items", e.Img), e)

		if slot, ok := h.lookupShipSlotByName(e.SlotName); ok {
			slot.PlaceItem(item)
		} else {
			h.dad.Slots[i+placed].PlaceItem(item)
			placed++
		}
	}

	h.recalculate()
}

func (h *hangar) lookupShipSlotByName(name string) (*fizzgui.DADSlot, bool) {
	for _, slot := range h.shipSlots {
		s := slot.UserData.(equip.Slot)
		if s.Name == name {
			return slot, true
		}
	}
	return nil, false
}

func (h *hangar) putOn(item *fizzgui.DADItem, slot *fizzgui.DADSlot, prevSlot *fizzgui.DADSlot) bool {
	if item.ID != slot.ID {
		//ID contains type of equipment, item possible to place in equal type, ex: weapon to weapon slot
		return false
	}

	slot.PlaceItem(item)
	h.recalculate()

	return true
}

func (h *hangar) takeOff(item *fizzgui.DADItem, slot *fizzgui.DADSlot, prevSlot *fizzgui.DADSlot) bool {

	slot.PlaceItem(item)
	h.recalculate()

	return true
}

func (h *hangar) recalculate() {
	p := player.Ship.InitParam.Param

	player.Ship.Equipment = nil
	player.Ship.LeftWeapon = nil
	player.Ship.RightWeapon = nil

	// h.ship.CurrParam = h.ship.InitParam

	for _, slot := range h.shipSlots {
		if slot.Item == nil {
			continue
		}

		switch slot.ID {
		case equip.TypeWeapon.Str():
			s := slot.UserData.(equip.Slot)

			w := slot.Item.UserData.(*weapons.Weapon)
			p.Weight += w.InitParam.Weight

			switch s.Side {
			case equip.Left:
				player.Ship.LeftWeapon = w
			case equip.Right:
				player.Ship.RightWeapon = w
			}

		default:
			e := slot.Item.UserData.(*equip.Equip)
			p = p.Summ(e.Param)

			player.Ship.Equipment = append(player.Ship.Equipment, e)
		}
	}

	h.infoTable.Class.update(player.Ship.Name, player.Ship.Class)
	h.infoTable.Weight.update(p.Weight)
	h.infoTable.Durability.update(p.Health)
	h.infoTable.Shield.update(p.Shield)
	h.infoTable.Speed.update(p.MovSpeed)
	h.infoTable.Contollability.update(p.RotSpeed)
	h.infoTable.Energy.update(p.Energy, p.EnergyAcc)
	h.infoTable.Metal.update(p.Metal, p.MetalAcc)

	h.infoTable.LeftWeapon.SetWeapon(equip.Left, player.Ship.LeftWeapon)
	h.infoTable.RightWeapon.SetWeapon(equip.Right, player.Ship.RightWeapon)
}

func (h *hangar) startBattle(wgt *fizzgui.Widget) {
	h.cMenu.Close()
	h.cShip.Close()
	h.cItems.Close()
	h.infoTable.C.Close()

	h.dad.Close()

	sceneIndex.Close()

	game.StartBattle(player)

	// Index()
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

	LeftWeapon  *infoWeapon
	RightWeapon *infoWeapon
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

	info.LeftWeapon = info.newInfoWeapon("Left")
	info.RightWeapon = info.newInfoWeapon("Right")

	return info
}

type infoWeapon struct {
	Name  *fizzgui.Widget
	DPS   *fizzgui.Widget
	Range *fizzgui.Widget
	Ammo  *fizzgui.Widget
}

func (info *informationTable) newInfoWeapon(side string) *infoWeapon {
	return &infoWeapon{
		Name:  info.newText(""),
		DPS:   info.newText(""),
		Range: info.newText(""),
		Ammo:  info.newText(""),
	}
}

func (iw *infoWeapon) SetWeapon(side equip.Side, w *weapons.Weapon) {
	if w == nil {
		iw.SetHidden(true)
		return
	}
	iw.SetHidden(false)

	iw.Name.Text = fmt.Sprintf("%s: %s", side, w.Name)
	iw.DPS.Text = fmt.Sprintf("  DPS: %.1f", w.InitParam.Damage/float32(w.InitParam.Rate.Seconds()))
	iw.Range.Text = fmt.Sprintf("  Range: %.1f", w.GetAttackRange(w.InitParam))
	iw.Ammo.Text = fmt.Sprintf("  Ammo: %d", w.InitParam.Ammunition)
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
