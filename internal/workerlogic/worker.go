package workerlogic

import (
	"log"
	"solune/tcp"
	"solune/store"
)

func Run(port string) {
	manager := store.NewDataStoreManager(string(port));

	log.Printf("Staring server on port %s\n", port);
	tcp.StartServer(port, manager);

	select {}
}
