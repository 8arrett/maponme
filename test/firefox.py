from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.wait import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.firefox.options import Options as FirefoxOptions
import sys

server = "http://localhost:8080"

def testTileSettings(browser, activeID, isMessagePresent):
    for tile in ["introTile", "permissionTile", "mapTile"]:
        displayed = browser.find_element(By.ID, tile).is_displayed()
        if tile == activeID:
            assert displayed, f"{tile} not found"
        else:
            assert not displayed, f"{tile} found"
    assert isMessagePresent == browser.find_element(By.ID, "messagePopup").is_displayed()

# # #  Create a shareable link with location override settings

ffOptions = FirefoxOptions()
ffOptions.set_preference('permissions.default.geo', 1) # flows past location permissions popup
ffOptions.set_preference('geo.provider.testing', True)
ffOptions.set_preference('geo.provider.network.url', 'data:application/json,{"location": {"lat": 40.7590, "lng": -73.9845}, "accuracy": 27000.0}')
host = webdriver.Firefox(options=ffOptions)

host.get(server)
assert host.current_url == server + "/"
testTileSettings(host, "introTile", False)

startButton = host.find_element(By.ID, "startButton")
assert startButton.is_displayed()
startButton.click()
testTileSettings(host, "mapTile", False)

pubID = host.execute_script("return host.id")
priKey = host.execute_script("return host.key")

# # #  A second browser running without location permissions

client = webdriver.Firefox()
client.get(server + "/#" + pubID)

testTileSettings(client, "mapTile", False)
assert client.find_element(By.ID, "messagePopup").text == ""

# # # Exit early and leave both windows running if flag set

if len(sys.argv) > 1 and sys.argv[1] == "--watch":
    print("All tests passed!")
    exit()

# # # Gracefully shut down shared link

closeButton = host.find_element(By.ID, "closeLink")
assert closeButton.is_displayed()
closeButton.click()

closeMsg = WebDriverWait(client, 5).until(EC.visibility_of_element_located((By.ID, "messagePopup"))).text
assert closeMsg[:36] == "  Your friend's map has disappeared!", closeMsg
testTileSettings(client, "introTile", True)
assert client.current_url == server + "/#"

testTileSettings(host, "introTile", False)
assert host.current_url == server + "/#"

# # # End

print("All tests passed!")
host.quit()
client.quit()