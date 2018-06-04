package example_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"gonum.org/v1/hdf5"
)

func ExampleLibVersion() {
	version, err := hdf5.LibVersion()
	if err != nil {
		fmt.Printf("** error ** %s\n", err)
		return
	}
	fmt.Printf("HDF5 version: %s", version)
}

func Example_extend_dims() {
	const (
		fname = "abc"
	)
	// create a file with a single 5x20 dataset
	f, err := hdf5.CreateFile(fname, hdf5.F_ACC_TRUNC)
	if err != nil {
		fmt.Println("uh oh")
	}
	defer f.Close()

	var data [100]uint16
	for i := range data {
		data[i] = uint16(i)
	}

	dims := []uint{20, 5}
	max_dims := []uint{40, 5}
	dspace, err := hdf5.CreateSimpleDataspace(dims, max_dims)
	if err != nil {
		fmt.Println("uh oh")
	}

	dset, err := f.CreateDataset("dset", hdf5.T_NATIVE_USHORT, dspace)
	if err != nil {
		fmt.Println("uh oh")
	}

	err = dset.Write(&data[0])
	if err != nil {
		fmt.Println("uh oh")
	}

	// newDims := []uint{40, 5}
	// dset.SetExtent(newDims)

	fmt.Println("n")
	//Output: n
}

func Example_write_read_structured_data() {
	const (
		fname  string = "SDScompound.h5"
		dsname string = "ArrayOfStructures"
		mbr1   string = "A_name"
		mbr2   string = "B_name"
		mbr3   string = "C_name"
		length uint   = 10
		rank   int    = 1
	)

	type s1Type struct {
		a int
		b float32
		c float64
		d [3]int
		e string
	}

	type s2Type struct {
		c float64
		a int
	}

	s1 := [length]s1Type{}
	for i := 0; i < int(length); i++ {
		s1[i] = s1Type{
			a: i,
			b: float32(i * i),
			c: 1. / (float64(i) + 1),
			d: [...]int{i, i * 2, i * 3},
			e: fmt.Sprintf("--%d--", i),
		}
	}

	// create data space
	dims := []uint{length}
	space, err := hdf5.CreateSimpleDataspace(dims, nil)
	if err != nil {
		panic(err)
	}

	// create the file
	f, err := hdf5.CreateFile(fname, hdf5.F_ACC_TRUNC)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Printf(":: file [%s] created (id=%d)\n", fname, f.Id())

	// create the memory data type
	dtype, err := hdf5.NewDatatypeFromValue(float32(1))
	if err != nil {
		panic("could not create a dtype")
	}

	// create the dataset
	dset, err := f.CreateDataset(dsname, dtype, space)
	if err != nil {
		panic(err)
	}
	fmt.Printf(":: dset (id=%d)\n", dset.Id())

	// write data to the dataset
	fmt.Printf(":: dset.Write...\n")
	err = dset.Write(&s1)
	if err != nil {
		panic(err)
	}
	fmt.Printf(":: dset.Write... [ok]\n")

	// release resources
	dset.Close()
	f.Close()

	// open the file and the dataset
	f, err = hdf5.OpenFile(fname, hdf5.F_ACC_RDONLY)
	if err != nil {
		panic(err)
	}
	dset, err = f.OpenDataset(dsname)
	if err != nil {
		panic(err)
	}

	// read it back into a new slice
	s2 := make([]float32, length)
	err = dset.Read(&s2)
	if err != nil {
		panic(err)
	}

	// display the fields
	fmt.Printf(":: data: %v\n", s2)

	// release resources
	dset.Close()
	f.Close()

	// Output:
	//   :: file [SDScompound.h5] created (id=72057594037927937)
	// :: dset (id=360287970189639681)
	// :: dset.Write...
	// :: dset.Write... [ok]
	// :: data: [0 0 0 0 0 1.875 0 0 0 0]
}

// make sure the examples have GODEBUG=cgocheck=0
func TestMain(m *testing.M) {
	if os.Getenv("GODEBUG") == "cgocheck=0" {
		os.Exit(m.Run())
	} else {
		fmt.Println("re-running `go test hdf5_example` with GODEBUG=cgocheck=0")
		os.Setenv("GODEBUG", "cgocheck=0")
		c := exec.Command("go", "test")
		out, err := c.Output()
		fmt.Println(string(out))
		if err != nil {
			os.Exit(1)
		}
	}
}
