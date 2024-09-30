package main

func main() {

	fileDecompressPtr := flag.Bool("decompress", false, "Should we decompress")
	fileCompressPtr := flag.Bool("compress", false, "should we compress")

	filePathInPtr := flag.String("in", "", "The File Path In")
	filePathOutPtr := flag.String("out", "", "The File Path Out")

	flag.Parse()

	// Example: ./main.exe -decompress -in=/path/to/input -out=/path/to/output
	if *fileDecompressPtr {
		if err := decompressZSTD(*filePathInPtr, *filePathOutPtr); err != nil {
			err = fmt.Errorf("failed to call decompress command: %v", err)
			panic(err)
		}
		// Example: ./main.exe -compress -in=/path/to/input -out=/path/to/output
	} else if *fileCompressPtr {
		if err := compressZSTD(*filePathInPtr, *filePathOutPtr); err != nil {
			err = fmt.Errorf("failed to call compress command: %v", err)
			panic(err)
		}
	} else {
		panic("You must specify a command")
	}
}

func compressZSTD(inPath, outPath string) error {

	fileIn, err := os.OpenFile(inPath, os.O_RDONLY, 0644)
	if err != nil {
		err = fmt.Errorf("you have an error reading file :D : %v", err)
		return err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(fileIn)

	fileOut, err := os.Create(outPath)
	if err != nil {
		err = fmt.Errorf("you have an error creating file: %v", err)
		return err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(fileOut)

	zstdWriter, err := zstd.NewWriter(fileOut)
	if err != nil {
		err = fmt.Errorf("you have an error writing to file: %v", err)
		return err
	}

	defer func(zstdWriter *zstd.Encoder) {
		_ = zstdWriter.Close()
	}(zstdWriter)

	if _, err = io.Copy(zstdWriter, fileIn); err != nil {
		err = fmt.Errorf("you have an error copying to file: %v", err)
		return err
	}

	return nil
}

func decompressZSTD(inPath, outPath string) error {

	fileIn, err := os.OpenFile(inPath, os.O_RDONLY, 0644)
	if err != nil {
		err = fmt.Errorf("you have an error reading file :D : %v", err)
		return err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(fileIn)

	fileOut, err := os.Create(outPath)
	if err != nil {
		err = fmt.Errorf("you have an error creating file: %v", err)
		return err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(fileOut)

	zstdReader, err := zstd.NewReader(fileIn, zstd.WithDecoderConcurrency(0))
	if err != nil {
		err = fmt.Errorf("you have an error reading file: %v", err)
		return err
	}

	defer zstdReader.Close()

	if _, err = io.Copy(fileOut, zstdReader); err != nil {
		err = fmt.Errorf("you have an error copying to file: %v", err)
		return err
	}

	return nil
}
