package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	rp, err := NewRepository("C:\\shortener_storage_test.txt")
	require.NoError(t, err)
	_, err = os.Stat("C:\\shortener_storage_test.txt")
	assert.NotEqual(t, os.ErrNotExist, err, "Файл для хранения сокращённых URL не существует")

	rp.Insert("hgfdstrjti345", "http://iste.ru")

	rp, err = NewRepository("C:\\shortener_storage_test.txt")
	require.NoError(t, err)

	val, err := rp.Select("hgfdstrjti345")
	require.NoError(t, err)
	assert.Equal(t, "http://iste.ru", val, "Не сохранены строки, записанные перед перезагрузкой")
}
