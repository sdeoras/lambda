package gan

import "github.com/sdeoras/api"

type Generator interface {
	Generate(count int) (*api.GanResponse, error)
}
