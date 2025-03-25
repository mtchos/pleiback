package handler

import (
	"encoding/json"
	"github.com/mtchos/pleiback/old/internal/adapter/integration/spotify"
	"github.com/mtchos/pleiback/old/internal/adapter/integration/spotify/mock"
	"github.com/mtchos/pleiback/old/internal/domain/entity"
	"github.com/mtchos/pleiback/old/internal/domain/service"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

var searchResponseStr = `
		[
			{
				"ID": "2nLtzopw4rPReszdYBJU6h",
				"Name": "Numb",
				"Artists": [
					{
						"ID": "6XyY86QOPPrYVGvF9ch6wz",
						"Name": "Linkin Park",
						"Genres": null,
						"URI": ""
					}
				],
				"URI": "client:track:2nLtzopw4rPReszdYBJU6h"
			},
			{
				"ID": "2HBBM75Xv3o2Mqdyh1NcM0",
				"Name": "Heavy Is the Crown",
				"Artists": [
					{
						"ID": "6XyY86QOPPrYVGvF9ch6wz",
						"Name": "Linkin Park",
						"Genres": null,
						"URI": ""
					}
				],
				"URI": "client:track:2HBBM75Xv3o2Mqdyh1NcM0"
			},
			{
				"ID": "57BrRMwf9LrcmuOsyGilwr",
				"Name": "Crawling",
				"Artists": [
					{
						"ID": "6XyY86QOPPrYVGvF9ch6wz",
						"Name": "Linkin Park",
						"Genres": null,
						"URI": ""
					}
				],
				"URI": "client:track:57BrRMwf9LrcmuOsyGilwr"
			}
		]`

func TestTrackHandler_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpotifyService := mock_spotify.mock_spotify.NewMockIntegration(ctrl)
	var searchResponse spotify.SearchTracksResponse
	err := json.Unmarshal([]byte(searchResponseStr), &searchResponse)
	if err != nil {
		t.Error("error unmarshalling json", "err", err)
	}

	mockSpotifyService.EXPECT().SearchTracks("Galinha Pintadinha", int64(10), int64(0)).Return(searchResponse, nil)

	trackService := service.NewTrackService(mockSpotifyService)
	trackHandler := NewTrackHandler(trackService)

	req := httptest.NewRequest("GET", "/tracks?q=Galinha%20Pintadinha", nil)
	w := httptest.NewRecorder()

	trackHandler.Search(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	var tracks []entity.Track
	err = json.NewDecoder(w.Result().Body).Decode(&tracks)
	if err != nil {
		t.Error("error decoding tracks", "err", err)
	}

	IDs := []string{"2nLtzopw4rPReszdYBJU6h", "2HBBM75Xv3o2Mqdyh1NcM0", "57BrRMwf9LrcmuOsyGilwr"}
	Names := []string{"Numb", "Heavy Is the Crown", "Crawling"}
	Artists := []string{"Linkin Park"}

	err = nil
	for _, track := range tracks {
		if !slices.Contains(IDs, track.ID) ||
			slices.Contains(Names, track.Name) ||
			slices.Contains(Artists, track.Artists[0].Name) {
			t.Errorf("some details are missing in track with ID %s", track.ID)
		}
	}
}
