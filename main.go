package main

import (
	"log"

	"github.com/jfreymuth/pulse"
	"gopkg.in/hraban/opus.v2"
)

func main() {
	enc, err := opus.NewEncoder(48000, 1, opus.AppVoIP)
	if err != nil {
		log.Panic(err)
	}
	if err := enc.SetMaxBandwidth(opus.SuperWideband); err != nil {
		log.Panic(err)
	}
	if err := enc.SetBitrate(32000); err != nil {
		log.Panic(err)
	}

	dec, err := opus.NewDecoder(48000, 1)
	if err != nil {
		log.Panic(err)
	}

	pa, err := pulse.NewClient()
	if err != nil {
		log.Panic(err)
	}
	defer pa.Close()

	bufPlay := &buffer{}
	play, err := pa.NewPlayback(
		func(p []int16) {
			bufPlay.Read(p)
		},
		pulse.PlaybackMono,
		pulse.PlaybackSampleRate(48000),
		pulse.PlaybackBufferSize(1920),
	)
	if err != nil {
		log.Panic(err)
	}
	play.Start()

	frame := make([]int16, 1920)
	data := make([]byte, 1920)
	frame2 := make([]int16, 1920)

	bufRec := &buffer{}
	rec, err := pa.NewRecord(
		func(p []int16) {
			bufRec.Write(p)
			if bufRec.Len() < 1920 {
				return
			}
			// has 20ms
			bufRec.Read(frame)
			nEnc, err := enc.Encode(frame, data)
			if err != nil {
				log.Panic(err)
			}

			nDec, err := dec.Decode(data[:nEnc], frame2)
			if err != nil {
				log.Panic(err)
			}
			bufPlay.Write(frame2[:nDec])

			// log.Printf("raw: %d, opus: %d, dec: %d", len(frame), nEnc, nDec)
		},
		pulse.RecordMono,
		pulse.RecordSampleRate(48000),
		pulse.RecordBufferFragmentSize(1920),
		pulse.RecordAdjustLatency(true),
	)
	if err != nil {
		log.Panic(err)
	}
	rec.Start()

	select {}
}
