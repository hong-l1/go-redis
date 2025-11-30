package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	rediscli "github.com/redis/go-redis/v9"
	"math"
	"os"
	"strconv"
)

func main() {
	ctx := context.Background()
	rdb := rediscli.NewClient(&rediscli.Options{
		Addr:          "localhost:6379", // Redis Stack 服务的地址
		UnstableResp3: true,
		Protocol:      2,
	})
	indexName := "eino_index"
	// 构建 KNN 查询
	// `*=>[KNN 2 @vector_content $blob]` 的含义:
	// - `*`: 匹配所有文档 (我们不过滤元数据)。
	// - `=>`: 表示这是一个混合查询，我们主要关心右边的向量部分。
	// - `[KNN 2 @vector_content $blob]`: 在 `vector_content` 字段上执行一个 K-最近邻查询，
	//   查找 2 个最近邻。`$blob` 是一个参数，我们将把查询向量的二进制数据传递给它。
	// DIALECT 2 是必须的，用于支持这种现代的查询语法。
	k := 2
	query := fmt.Sprintf("*=>[KNN %d @content_vector $blob AS score]", k)
	searchContent := "golang goole开源 "
	embedder, err := ark.NewEmbedder(context.Background(), &ark.EmbeddingConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  "doubao-embedding-large-text-250515",
	})
	if err != nil {
		panic(err)
	}
	embeddings, err := embedder.EmbedStrings(ctx, []string{searchContent})
	if err != nil {
		panic(err)
	}
	searchResult, err := rdb.FTSearchWithArgs(ctx, indexName, query, &rediscli.FTSearchOptions{
		Params: map[string]interface{}{
			"blob": vector2Bytes(embeddings[0]),
		},
		DialectVersion: 2,
		Return: []rediscli.FTSearchReturn{
			{
				FieldName: "content",
			},
			{
				FieldName: "score",
			},
		},
	}).Result()
	if err != nil {
		panic(err)
	}
	for _, v := range searchResult.Docs {
		dist, _ := strconv.ParseFloat(v.Fields["score"], 64)
		sim := 1 - dist
		fmt.Printf("内容: %v | 距离: %.6f | 相似度: %.6f\n", v.Fields["content"], dist, sim)
	}

}
func vector2Bytes(vector []float64) []byte {
	float32Arr := make([]float32, len(vector))
	for i, v := range vector {
		float32Arr[i] = float32(v)
	}
	bytes := make([]byte, len(float32Arr)*4)
	for i, v := range float32Arr {
		binary.LittleEndian.PutUint32(bytes[i*4:], math.Float32bits(v))
	}
	return bytes
}
