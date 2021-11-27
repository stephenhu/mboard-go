package main

import (
	"errors"
	"fmt"
  "io"
	"log"
	"os/exec"

	"github.com/madsportslab/glbs"
)

const (
	OMX_PLAYER						= "omxplayer"
	OMX_HDMI							= "hdmi"
	OMX_OPTION_O  				= "-o"
	OMX_OPTION_D  				= "-d"
	OMX_DBUS_DIR  				= "/org/mpris/MediaPlayer2"
	OMX_DBUS_PLAYER				= "org.mpris.MediaPlayer2.Player.Action"
  OMX_DBUS_STOP     		= 15
)

var (
  omx *exec.Cmd
  mctl io.WriteCloser
)


func videoPlay(options map[string]string) {

	if omx != nil {
		log.Println("Video already playing")
	}

	glbs.SetNamespace("blobs")

	filename := fmt.Sprintf("/home/mboard/bin/%s", *glbs.GetPath(options["key"]))

	log.Println(filename)

	omx = exec.Command(OMX_PLAYER, OMX_OPTION_O, OMX_HDMI, filename)

	ctl, err := omx.StdinPipe()

  if err != nil {
    log.Println(err)
    return
  }

  err = omx.Start()

	if err != nil {
    log.Println(err)
		return
	}

  mctl = ctl

	err = omx.Wait()

	if err != nil {
    log.Println(err)
		return
	}

} // videoPlay

func videoStop() error {

  log.Println("stopping video")

	if omx == nil {
		return errors.New("No video playing")
	}

  _, err := io.WriteString(mctl, "q")

	if err != nil {
		return err
	}

	return nil

} // videoStop
