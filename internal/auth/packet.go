package auth

import (
	"flase_api/internal/db"
	"flase_api/internal/encoder"
	"flase_api/internal/packets"
	smartscheduler "flase_api/internal/smart_scheduler"
	"log"
)

type PacketSystem struct {
	loadedPackets map[string]*packets.Packet
}

func NewPacketSystem() *PacketSystem {
	return &PacketSystem{
		loadedPackets: make(map[string]*packets.Packet),
	}
}

func (usr *User) LoadPacket(ctx *db.DbCtx, pass string, sched *smartscheduler.SmartScheduler) (*packets.Packet, error) {
	hPass := encoder.Enc(pass)
	if pac, exists := usr.packetSystem.loadedPackets[hPass]; exists {
		if pac.IsDestructed() {
			delete(usr.packetSystem.loadedPackets, hPass)
		} else {
			return pac, nil
		}
	}
	log.Println("trying to load packet: ", ctx.GetChild(hPass).GetPath())
	pac, err := packets.DbLoad(ctx.GetChild(hPass), sched)
	if err != nil {
		return nil, err
	}
	usr.packet = pac
	usr.packetSystem.loadedPackets[hPass] = pac
	return pac, nil
}

func (usr *User) GetPacket() *packets.Packet {
	return usr.packet
}
