import aiohttp

API_URL = 'http://localhost:3000'
API_CATEGORIES = '/category'
API_ANNOUNCEMENT = '/announcement'
API_ANNOUNCEMENT_UPDATE = API_ANNOUNCEMENT + '/UpdatePrices'
API_QUERY = '/query'
API_IDENTIFICATION = '/identification'

async def fetch(session, url, params = {}):
    async with session.get(url, params=params) as response:
        json = await response.json()
        try:
            return json['result'], False
        except:
            return json, True

async def getAllAnnouncements(session):
    url = API_URL + API_ANNOUNCEMENT
    return await fetch(session, url)

async def getAnnouncementsByCategory(session, categoryName):
    url = API_URL + API_ANNOUNCEMENT
    params = {'category': categoryName}
    return await fetch(session, url, params)

async def getAnnouncementsByCategoryLt(session, categoryName, maxPrice):
    url = API_URL + API_ANNOUNCEMENT
    params = {'category': categoryName, 'lt': str(maxPrice)}
    return await fetch(session, url, params)

async def getCategories(session):
    url = API_URL + API_CATEGORIES
    return await fetch(session, url)

async def getCategoryByName(session, name):
    url = API_URL + API_CATEGORIES
    params = {'categoryName': name}
    return await fetch(session, url, params)

async def getIdentification(session):
    url = API_URL + API_IDENTIFICATION
    async with session.get(url) as response:
        return await response.json()

async def getQueriesByAnnouncement(session, announcementId):
    url = API_URL + API_QUERY
    params = {'announcementId': announcementId}
    return await fetch(session, url, params)

async def getQueriesByIssuer(session, id):
    url = API_URL + API_QUERY
    params = {'issuerId': id}
    return await fetch(session, url, params)

async def getQueryById(session, queryId):
    url = API_URL + API_QUERY
    params = {'queryId': queryId}
    return await fetch(session, url, params)

async def makeQuery(session, announcementId, query, price):
    url = API_URL + API_QUERY
    jsonData = {
        'announcementId': announcementId,
        'query': query,
        'price': price
    }
    result = await session.post(url, json=jsonData)
    responseJson = await result.json()
    try:
        return responseJson['result']
    except:
        return responseJson    

async def makeAnnouncement(session, filePath, filename, categoryName, queries):
    url = API_URL + API_ANNOUNCEMENT
    data = aiohttp.FormData()
    data.add_field('data_file',
               open(filePath, 'rb'),
               filename=filename)
    data.add_field('queries', str(queries).replace('\'','"'))
    data.add_field('category', categoryName)
    result = await session.post(url, data=data)
    responseJson = await result.json()
    try:
        return responseJson['result']
    except:
        return responseJson

async def updatePrices(session, announcementId, updates):
    url = API_URL + API_ANNOUNCEMENT_UPDATE
    jsonData = {
        'announcementId': announcementId,
        'updates': updates,
    }
    result = await session.post(url, json=jsonData)
    responseJson = await result.json()
    try:
        return responseJson['result']
    except:
        return responseJson  




        
