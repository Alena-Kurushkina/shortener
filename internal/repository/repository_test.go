package repository

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	rp, err := newFileRepository("/Users/alena/shortener_storage_test.txt")
	require.NoError(t, err)
	_, err = os.Stat("/Users/alena/shortener_storage_test.txt")
	assert.NotEqual(t, os.ErrNotExist, err, "Файл для хранения сокращённых URL не существует")

	rp.Insert(context.TODO(),"hgfdstrjti345", "http://iste.ru")

	rp, err = newFileRepository("/Users/alena/shortener_storage_test.txt")
	require.NoError(t, err)

	val, err := rp.Select(context.TODO(),"hgfdstrjti345")
	require.NoError(t, err)
	assert.Equal(t, "http://iste.ru", val, "Не сохранены строки, записанные перед перезагрузкой")
}
