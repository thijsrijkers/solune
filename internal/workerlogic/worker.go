package workerlogic

import (
	"log"
	"solune/store"
	"solune/tcp"
)

func Run(port string) {
	manager := store.NewDataStoreManager()

	log.Printf("Staring server on port %s\n", port)
	tcp.StartServer(port, manager)

	select {}
}
