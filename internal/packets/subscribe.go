package packets

import (
	messagereadwrite "flase_api/internal/message_read_write"
	"fmt"
	"log"
)

type packetInfo struct {
	Title        string `json:"title"`
	DeliveryTime string `json:"deliveryTime"`
}

func (pac *Packet) PacketInfo() *packetInfo {
	return &packetInfo{
		Title:        pac.name,
		DeliveryTime: pac.deliveryTime.Format(TimeLayout),
	}
}

func (pac *Packet) onUpdated() {
	if len(pac.subscribers) == 0 {
		return
	}
	info := pac.PacketInfo()

	log.Println(pac.subscribers)

	toDelete := make([]string, 0)
	for sessIds, rw := range pac.subscribers {
		if rw.IsClosed() {
			toDelete = append(toDelete, sessIds)
			continue
		}
		rw.WriteAny("packet", info)
	}

	for _, ssid := range toDelete {
		delete(pac.subscribers, ssid)
	}
}

func (pac *Packet) Subscribe(sessId string, rw *messagereadwrite.ReadWriter) {
	fmt.Println("subscribe")
	pac.subscribers[sessId] = rw
	info := pac.PacketInfo()
	rw.WriteAny("packet", info)
}

func (pac *Packet) Unsubscribe(sessId string) {
	delete(pac.subscribers, sessId)
}
