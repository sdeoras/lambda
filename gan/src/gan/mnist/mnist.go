package mnist

import (
	"gan/src/gan"
	"io/ioutil"

	"github.com/gogo/protobuf/proto"
	"github.com/sdeoras/api/pb"
	ganop "github.com/sdeoras/comp/gan"
	"github.com/sdeoras/comp/gan/ganMnist"
)

type generator struct {
	model string
	op    ganop.Operator
}

// NewGenerator returns a new instance of GAN generator interface.
// It wraps ganop Operator and loads it with checkpoint model.
func NewGenerator(model string) (gan.Generator, error) {
	g := new(generator)
	op, err := ganMnist.NewOperator(nil)
	if err != nil {
		return nil, err
	}
	g.op = op

	g.model = model

	b, err := ioutil.ReadFile(model)
	if err != nil {
		return nil, err
	}

	cp := new(pb.Checkpoint)
	if err := proto.Unmarshal(b, cp); err != nil {
		return nil, err
	}

	if err := g.op.Load(cp); err != nil {
		return nil, err
	}

	return g, nil
}

// Generate implements GAN generator interface
func (g *generator) Generate(count int) (*pb.GanResponse, error) {
	bs, err := g.op.Generate(count)
	if err != nil {
		return nil, err
	}

	out := new(pb.GanResponse)
	out.Images = make([]*pb.Image, len(bs))

	for i := range bs {
		out.Images[i] = new(pb.Image)
		out.Images[i].Data = bs[i]
	}

	return out, nil
}
