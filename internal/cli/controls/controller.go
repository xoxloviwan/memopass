package controls

import (
	"encoding/json"
	iCrypto "iwakho/gopherkeep/internal/cli/crypto"
	iHttp "iwakho/gopherkeep/internal/cli/http"
	"iwakho/gopherkeep/internal/model"
	"os"
)

type Controller struct {
	cli    *iHttp.Client
	crypto *iCrypto.CryptoManager
}

func New(cli *iHttp.Client) *Controller {
	return &Controller{cli, nil}
}

func (c *Controller) Login(p model.Pair) error {
	err := c.cli.Login(p)
	if err != nil {
		return err
	}
	c.crypto = iCrypto.NewCryptoManager(p)
	return nil
}

func (c *Controller) SignUp(p model.Pair) error {
	err := c.cli.SignUp(p)
	if err != nil {
		return err
	}
	c.crypto = iCrypto.NewCryptoManager(p)
	return nil
}

func (c *Controller) AddCard(card model.Card) error {
	body, header, err := model.FillCardForm(card, c.crypto)
	if err != nil {
		return err
	}
	return c.cli.AddItem(c.cli.Add.Card, body, header)
}

func (c *Controller) AddPair(p model.Pair) error {
	body, header, err := model.FillPairForm(p, c.crypto)
	if err != nil {
		return err
	}
	return c.cli.AddItem(c.cli.Add.Pair, body, header)
}

func (c *Controller) AddFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	body, header, err := model.FillFileForm(file)
	if err != nil {
		return err
	}
	return c.cli.AddItem(c.cli.Add.File, body, header)
}

func (c *Controller) GetPairs(limit int, offset int) ([]model.PairInfo, error) {
	pairs, err := c.cli.GetItems(c.cli.Api.Get.Pair, limit, offset)
	if err != nil {
		return nil, err
	}
	return model.DecryptPairs(pairs, c.crypto)
}

func (c *Controller) GetCards(limit int, offset int) ([]model.CardInfo, error) {
	cards, err := c.cli.GetItems(c.cli.Api.Get.Card, limit, offset)
	if err != nil {
		return nil, err
	}
	return model.DecryptCards(cards, c.crypto)
}

func (c *Controller) AddText(text string) error {
	body, header, err := model.FillTextForm(text)
	if err != nil {
		return err
	}
	return c.cli.AddItem(c.cli.Add.Text, body, header)
}

func (c *Controller) GetTexts(limit int, offset int) ([]model.FileInfo, error) {
	data, err := c.cli.GetItems(c.cli.Api.Get.Text, limit, offset)
	if err != nil {
		return nil, err
	}
	var files []model.FileInfo
	err = json.Unmarshal(data, &files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (c *Controller) GetFiles(limit int, offset int) ([]model.FileInfo, error) {
	data, err := c.cli.GetItems(c.cli.Api.Get.File, limit, offset)
	if err != nil {
		return nil, err
	}
	var files []model.FileInfo
	err = json.Unmarshal(data, &files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (c *Controller) GetFileById(id int) (*model.File, error) {
	file := new(model.File)
	var err error
	file.Blob, file.Name, err = c.cli.GetFileById(c.cli.Api.GetById.File, id)
	return file, err
}

func (c *Controller) GetTextById(id int) (*model.File, error) {
	file := new(model.File)
	var err error
	file.Blob, file.Name, err = c.cli.GetFileById(c.cli.Api.GetById.Text, id)
	return file, err
}
