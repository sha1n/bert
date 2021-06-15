package pkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_PerceivedTimeCPUTimer_Unloaded(t *testing.T) {
	timer := PerceivedTimeCPUTimer{}

	elapsed := timer.Start()
	perceived, usr, sys := elapsed()

	assert.GreaterOrEqual(t, time.Millisecond*1, perceived)
	assert.Equal(t, time.Nanosecond*0, usr)
	assert.Equal(t, time.Nanosecond*0, sys)
}

func Test_PerceivedTimeCPUTimer_Loaded(t *testing.T) {
	timer := PerceivedTimeCPUTimer{}

	elapsed := timer.Start()

	assert.Eventually(t,
		func() bool {
			perceived, usr, sys := elapsed()

			return perceived > time.Millisecond*1 && usr == time.Nanosecond*0 && sys == time.Nanosecond*0
		},
		time.Second*1,
		time.Microsecond*1,
	)
}
