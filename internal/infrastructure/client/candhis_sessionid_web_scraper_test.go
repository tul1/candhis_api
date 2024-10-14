package client_test

// func TestGetCandhisSessionID_Success(t *testing.T) {
// 	chromeURL := "localhost:9222"
// 	chromeID := "browser-id"
// 	targetWeb := "https://candhis.cerema.fr"
// 	expectedSessionID := "test-php-session-id"

// 	patch := monkey.Patch(chromedp.Run, func(ctx context.Context, actions ...chromedp.Action) error {
// 		for i := range actions {
// 			if i == 2 {
// 				cookies := []*network.Cookie{
// 					{Name: "PHPSESSID", Value: expectedSessionID},
// 				}
// 				ctx = context.WithValue(ctx, "cookies", cookies)
// 			}
// 		}
// 		return nil
// 	})
// 	defer patch.Unpatch()

// 	candhisClient := client.NewCandhisSessionIDWebScraper(chromeURL, chromeID, targetWeb)

// 	sessionID, err := candhisClient.GetCandhisSessionID(context.Background())

// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedSessionID, sessionID.ID())
// }

// func TestGetCandhisSessionID_Failure_NoCookies(t *testing.T) {
// 	// Arrange
// 	chromeURL := "localhost:9222"
// 	chromeID := "browser-id"
// 	targetWeb := "https://candhis.cerema.fr"

// 	// Patch chromedp.Run to simulate cookie retrieval failure
// 	patch := monkey.Patch(chromedp.Run, func(ctx context.Context, actions ...chromedp.Action) error {
// 		for _, action := range actions {
// 			switch act := action.(type) {
// 			case chromedp.ActionFunc:
// 				// Simulate failure in ActionFunc that retrieves cookies
// 				return nil
// 			}
// 		}
// 		return nil
// 	})
// 	defer patch.Unpatch()

// 	// Patch the Do() method on network.GetCookiesParams to simulate no cookies being returned
// 	patchGetCookies := monkey.PatchInstanceMethod(
// 		reflect.TypeOf(&network.GetCookiesParams{}),
// 		"Do",
// 		func(_ *network.GetCookiesParams, ctx context.Context) ([]*network.Cookie, error) {
// 			return []*network.Cookie{}, nil
// 		},
// 	)
// 	defer patchGetCookies.Unpatch()

// 	// Create the actual client
// 	candhisClient := client.NewCandhisSessionIDWebScraper(chromeURL, chromeID, targetWeb)

// 	// Act
// 	sessionID, err := candhisClient.GetCandhisSessionID(context.Background())

// 	// Assert
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "failed to retrieve session id")
// 	assert.Empty(t, sessionID.ID)
// }

// func TestGetCandhisSessionID_Failure_ChromedpError(t *testing.T) {
// 	// Arrange
// 	chromeURL := "localhost:9222"
// 	chromeID := "browser-id"
// 	targetWeb := "https://candhis.cerema.fr"

// 	// Patch chromedp.Run to simulate a failure
// 	patch := monkey.Patch(chromedp.Run, func(ctx context.Context, actions ...chromedp.Action) error {
// 		return errors.New("chromedp failed")
// 	})
// 	defer patch.Unpatch()

// 	// Create the actual client
// 	candhisClient := client.NewCandhisSessionIDWebScraper(chromeURL, chromeID, targetWeb)

// 	// Act
// 	sessionID, err := candhisClient.GetCandhisSessionID(context.Background())

// 	// Assert
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "chromedp failed")
// 	assert.Empty(t, sessionID.ID)
// }
