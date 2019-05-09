package mardown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_numbering(t *testing.T) {
	var n numbering

	assert.Equal(t, "", n.Render())
	n.NextLevel(1)
	assert.Equal(t, "1", n.Render())
	n.NextLevel(1)
	assert.Equal(t, "2", n.Render())
	n.NextLevel(1)
	assert.Equal(t, "3", n.Render())
	n.NextLevel(2)
	assert.Equal(t, "3.1", n.Render())
	n.NextLevel(2)
	assert.Equal(t, "3.2", n.Render())
	n.NextLevel(2)
	assert.Equal(t, "3.3", n.Render())
	n.NextLevel(4)
	assert.Equal(t, "3.3.0.1", n.Render())
	n.NextLevel(4)
	assert.Equal(t, "3.3.0.2", n.Render())
	n.NextLevel(3)
	assert.Equal(t, "3.3.1", n.Render())
	n.NextLevel(3)
	assert.Equal(t, "3.3.2", n.Render())
	n.NextLevel(1)
	assert.Equal(t, "4", n.Render())
}
