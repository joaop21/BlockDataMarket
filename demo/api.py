import aiohttp

API_CATEGORIES = '/category'
API_ANNOUNCEMENT = '/announcement'
API_ANNOUNCEMENT_UPDATE = API_ANNOUNCEMENT + '/UpdatePrices'
API_QUERY = '/query'
API_IDENTIFICATION = '/identification'

class API:

    def __init__(self, api_url, session):
        self.API_URL = api_url
        self.session = session
        self.queries = []

    async def fetch(self, url, params={}):
        async with self.session.get(url, params=params) as response:
            json = await response.json()
            try:
                return json['result'], False
            except:
                return json, True

    async def post(self, url, jsonData):
        result = await self.session.post(url, json=jsonData)
        responseJson = await result.json()
        try:
            return responseJson['result']
        except:
            return responseJson


    async def getAllAnnouncements(self):
        url = self.API_URL + API_ANNOUNCEMENT
        return await self.fetch(url)

    async def getAnnouncementsByCategory(self, categoryName):
        url = self.API_URL + API_ANNOUNCEMENT
        params = {'category': categoryName}
        return await self.fetch(url, params)

    async def getAnnouncementsByCategoryLt(self, categoryName, maxPrice):
        url = self.API_URL + API_ANNOUNCEMENT
        params = {'category': categoryName, 'lt': str(maxPrice)}
        return await self.fetch(url, params)

    async def getCategories(self):
        url = self.API_URL + API_CATEGORIES
        return await self.fetch(url)

    async def getCategoryByName(self, name):
        url = self.API_URL + API_CATEGORIES
        params = {'categoryName': name}
        return await self.fetch(url, params)

    async def getIdentification(self):
        url = self.API_URL + API_IDENTIFICATION
        async with self.session.get(url) as response:
            return await response.json()

    async def getQueriesByAnnouncement(self, announcementId):
        url = self.API_URL + API_QUERY
        params = {'announcementId': announcementId}
        return await self.fetch(url, params)

    async def getQueriesByIssuer(self, id):
        url = self.API_URL + API_QUERY
        params = {'issuerId': id}
        return await self.fetch(url, params)

    async def getQueryById(self, queryId):
        url = self.API_URL + API_QUERY
        params = {'queryId': queryId}
        return await self.fetch(url, params)

    async def makeQuery(self, announcementId, query, price):
        url = self.API_URL + API_QUERY
        jsonData = {
            'announcementId': announcementId,
            'query': query,
            'price': price
        }
        return await self.post(url, jsonData)


    async def makeAnnouncement(self, filePath, filename, categoryName, queries):
        url = self.API_URL + API_ANNOUNCEMENT
        data = aiohttp.FormData()
        data.add_field('data_file',
                       open(filePath, 'rb'),
                       filename=filename)
        data.add_field('queries', str(queries).replace('\'', '"'))
        data.add_field('category', categoryName)
        result = await self.session.post(url, data=data)
        responseJson = await result.json()
        try:
            return responseJson['result']
        except:
            return responseJson

    async def updatePrices(self, announcementId, updates):
        url = self.API_URL + API_ANNOUNCEMENT_UPDATE
        jsonData = {
            'announcementId': announcementId,
            'updates': updates,
        }
        return await self.post(url, jsonData)

    async def makeIdentification(self, name):
        url = self.API_URL + API_IDENTIFICATION
        jsonData = {
            'name': name,
        }
        return await self.post(url, jsonData)

    async def makeCategory(self, name, queries):
        url = self.API_URL + API_CATEGORIES
        jsonData = {
            'name': name,
            'queries': str(queries).replace('\'', '"')
        }
        return await self.post(url, jsonData)           

