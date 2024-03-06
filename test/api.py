import requests

# # # $ pytest api.py

server = 'http://localhost:8080'
jsonHeader = {'content-type': 'application/json'}

### The initial request should always return temporary credentials if the server is running

postInit = requests.post(server + '/api/', headers=jsonHeader)

def test_postInit():
    assert postInit.status_code == 200
    assert len(postInit.json()['id']) == 8
    assert len(postInit.json()['key']) == 24

### Create keys for subsequent requests

targetID = postInit.json()['id']
targetKey = postInit.json()['key']

badKey = targetKey[23]
if badKey == '+' or badKey == '/' or badKey == 'z' or badKey == 'Z':
    badKey = 'a'
badKey = targetKey[:23] + chr(ord(badKey) + 1)

### Once these credentials open a public link, it should return a null location

getEmpty = requests.get(server + '/api/' + targetID, headers=jsonHeader)

def test_getEmpty():
    assert getEmpty.status_code == 200
    assert getEmpty.json()['lon'] == ''
    assert getEmpty.json()['lat'] == ''
    assert getEmpty.json()['acc'] == ''
    assert getEmpty.json()['mod'] >= 0

### Once these credentials open a public link, it should not allow changing of location position from a bad key

badLocation = {'Key': badKey, 'Lat': '12.34', 'Lon': '23.45', 'Acc': '100'}
putFirstBad = requests.put(server + '/api/' + targetID, headers=jsonHeader, json=badLocation)

def test_putFirstBad():
    assert putFirstBad.status_code == 400
    assert putFirstBad.json()['error'] == 'data-auth'

### Once these credentials open a public link, it should allow changing of location position

firstLocation = {'Key': targetKey, 'Lat': '12.34', 'Lon': '23.45', 'Acc': '100'}
putFirst = requests.put(server + '/api/' + targetID, headers=jsonHeader, json=firstLocation)

def test_putFirst():
    assert putFirst.status_code == 200
    assert putFirst.json()['mod'] >= 1600000000

### Verify that the previous change was correctly set

getFirst = requests.get(server + '/api/' + targetID, headers=jsonHeader)

def test_getFirst():
    assert getFirst.status_code == 200
    assert getFirst.json()['lon'] == firstLocation['Lon']
    assert getFirst.json()['lat'] == firstLocation['Lat']
    assert getFirst.json()['acc'] == firstLocation['Acc']
    assert getFirst.json()['mod'] >= 1600000000

### Once a location changes, it should not allow continued changes for a moment

secondLocation = {'Key': targetKey, 'Lat': '98.76', 'Lon': '87.65', 'Acc': '200'}
putSecond = requests.put(server + '/api/' + targetID, headers=jsonHeader, json=secondLocation)

def test_putSecond():
    assert putSecond.status_code == 400
    assert putSecond.json()['error'] == 'data-cooldown'

### When a quick second location change is sent, a viewer should only see the original change

getSecond = requests.get(server + '/api/' + targetID, headers=jsonHeader)

def test_getSecond():
    assert getSecond.status_code == 200
    assert getSecond.json()['lon'] == firstLocation['Lon']
    assert getSecond.json()['lat'] == firstLocation['Lat']
    assert getSecond.json()['acc'] == firstLocation['Acc']
    assert getSecond.json()['mod'] >= 1600000000

### When a request to close a credential comes from a valid key, it should succeed

finalAuth = {'Key': targetKey}
deleteFinal = requests.delete(server + '/api/' + targetID, headers=jsonHeader, json=finalAuth)

def test_deleteFinal():
    assert deleteFinal.status_code == 200
    assert deleteFinal.json()['status'] == 'success'

### When a credential has been removed, it should no longer respond to Get requests

getFinal = requests.get(server + '/api/' + targetID, headers=jsonHeader)

def test_getFinal():
    assert getFinal.status_code == 400
    assert getFinal.json()['error'] == 'data-expired'

### When a credential has been removed, it should no longer respond to Put requests

putFinal = requests.put(server + '/api/' + targetID, headers=jsonHeader, json=firstLocation)

def test_putFinal():
    assert putFinal.status_code == 400
    assert putFinal.json()['error'] == 'data-expired'

### When a credential has been removed, it should no longer respond to Delete requests

deleteAgain = requests.delete(server + '/api/' + targetID, headers=jsonHeader, json=finalAuth)

def test_putFinal():
    assert putFinal.status_code == 400
    assert putFinal.json()['error'] == 'data-expired'