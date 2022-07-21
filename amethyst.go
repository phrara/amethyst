package amethyst

import (
	"github.com/phrara/amethyst/config"
	"strconv"
)

const (
	LOGO = ` 
   _____                  __  .__                    __   
  /  _  \   _____   _____/  |_|  |__ ___.__. _______/  |_ 
 /  /_\  \ /     \_/ __ \   __\  |  <   |  |/  ___/\   __\
/    |    \  Y Y  \  ___/|  | |   Y  \___  |\___ \  |  |  
\____|__  /__|_|  /\___  >__| |___|  / ____/____  > |__|  
        \/      \/     \/          \/\/         \/        
`
)

func New() *Server {
	return generate(LOGO, config.Global.IP, strconv.Itoa(config.Global.Port), config.Global.Protocol)
}
