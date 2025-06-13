package service

import (
	"container/heap"
	"fmt"
	"tfidf/internal/model"
)

func HuffmanEncode(document []byte) (string, *model.Node, error) {
	if len(document) == 0 {
		return "", nil, fmt.Errorf("отсутствует содержимое документа")
	}

	freqMap := countWordFrequency(document)

	root, err := makeHuffmanTree(freqMap)
	if err != nil {
		return "", nil, err
	}

	codeMap := make(map[byte]string)
	generateCodes(root, "", codeMap)

	code := ""
	for _, char := range document {
		code += codeMap[char]
	}

	return code, root, nil
}

func countWordFrequency(content []byte) map[byte]int {
	freq := make(map[byte]int)
	for _, b := range content {
		freq[b]++
	}
	return freq
}

func makeHuffmanTree(freqMap map[byte]int) (*model.Node, error) {
	h := model.NewMinHeapWithFreqMap(freqMap)

	for h.Len() > 1 {
		l := heap.Pop(h).(*model.Node)
		r := heap.Pop(h).(*model.Node)

		merged := &model.Node{
			Freq:  l.Freq + r.Freq,
			Left:  l,
			Right: r,
		}
		heap.Push(h, merged)
	}

	root := heap.Pop(h).(*model.Node)
	return root, nil
}

func generateCodes(node *model.Node, prefix string, codeMap map[byte]string) {
	if node.Left == nil && node.Right == nil {
		codeMap[node.Char] = prefix
		return
	}

	if node.Left != nil {
		generateCodes(node.Left, prefix+"0", codeMap)
	}

	if node.Right != nil {
		generateCodes(node.Right, prefix+"1", codeMap)
	}
}
