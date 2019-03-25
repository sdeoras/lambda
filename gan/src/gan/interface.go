package gan

import (
	"github.com/sdeoras/api/pb"
)

type Generator interface {
	Generate(count int) (*pb.GanResponse, error)
}
