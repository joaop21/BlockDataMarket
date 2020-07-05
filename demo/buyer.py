import aiohttp
import asyncio
import pprint
import random
from api import API
import sys
import argparse
from datetime import datetime

pp = pprint.PrettyPrinter(indent=2)

class Buyer:

    def __init__(self, id, api):
        self.api = api
        self.category = {}
        self.id = id
        self.queries = []

    async def pickCategory(self):
        categories, _ = await self.api.getCategories()
        index = random.randint(0, len(categories)-1)
        self.category = categories[index]

    async def cheaperQuery(self):
        query = random.choice(self.category['possibleQueries'])
        announcements = await self.api.getAnnouncementsByCategory(self.category['name'])
        announcementId, price = Buyer.selectCheaperAnnouncement(announcements, query)
        resultQuery = await self.api.makeQuery(announcementId, query, price)
        self.queries.append(resultQuery['queryId'])
        printOp("CheaperQuery", self.id, announcementId, query, price)

    async def randomQuery(self):
        announcements, _ = await self.api.getAllAnnouncements()
        announcement = random.choice(announcements)
        index = random.randrange(len(announcement['possibleQueries']))
        query = announcement['possibleQueries'][index]
        queryPrice = announcement['prices'][index]
        priceToPay = random.randint(round(queryPrice/2),queryPrice)
        resultQuery = await self.api.makeQuery(announcement['announcementId'], query, priceToPay)
        self.queries.append(resultQuery['queryId'])
        printOp("RandomQuery", self.id, announcement['announcementId'], query, priceToPay)

    async def queryLowerThan(self):
        announcements, _ = await self.api.getAllAnnouncements()
        query = random.choice(self.category['possibleQueries'])
        meanPrice = await Buyer.getQueryMeanPrice(announcements, query)
        announcements, _ = await self.api.getAnnouncementsByCategoryLt(self.category['name'], meanPrice)
        announcement = random.choice(announcements)
        query, price = Buyer.getCheapestQuery(announcement)
        resultQuery = await self.api.makeQuery(announcement['announcementId'], query, price)
        self.queries.append(resultQuery['queryId'])
        printOp("QueryLowerThan", self.id, announcement['announcementId'], query, price)

    async def consultPreviousQuery(self):
        try:
            queryId = random.choice(self.queries)
            query, _ = await self.api.getQueryById(queryId)
            printOp("ConsultPreviousQuery", self.id)
        except:
            pass    

    @staticmethod
    def selectCheaperAnnouncement(announcements, query):
        prices = []
        for announcement in announcements:
            for tuple in zip(announcement['possibleQueries'], announcement['prices']):
                if tuple[0] == query:
                    prices.append((announcement['announcementId'], tuple[1]))
        prices.sort(key=lambda x:x[1])
        return prices[0][0], prices[0][1]

    @staticmethod
    async def getQueryMeanPrice(announcements, query):
        prices = []
        for announcement in announcements:
            try:
                index = announcement['possibleQueries'].index(query)
                prices.append(announcement['prices'][index])
            except:
                pass
        return sum(prices) / len(prices)

    @staticmethod
    def getCheapestQuery(announcement):
        prices = zip(announcement['possibleQueries'], announcement['prices'])
        prices = sorted(prices, key = lambda t: t[1])
        return prices[0][0], prices[0][1]
        
    async def lifeCycle(self, numOps):
        if not self.category:
            await self.pickCategory()

        for _ in range(numOps):
            try:
                await asyncio.sleep(5)
                op = random.random()
                if op < 0.4:
                    await self.cheaperQuery()
                elif op < 0.7:
                    await self.randomQuery()
                elif op < 0.9:
                    await self.queryLowerThan()
                else:
                    await self.consultPreviousQuery()
            except:
                pass                 

def printOp(op, id, announcementId=-1, query="", price=-1):
    now = datetime.now()
    current_time = now.strftime("%H:%M:%S")
    if(price != -1):
        print(current_time + " -> (" + op +") Buyer " + str(id) + " queried " + query + " on announcement " + announcementId + " for " + str(price))
    else:
        print(current_time + " -> (" + op +") Buyer " + str(id) + " consulted previous query")


            

async def main(api_url, buyers, actions):
    async with aiohttp.ClientSession() as session:
        api = API(api_url, session)
        tasks = []
        for i in range(buyers):
            buyer = Buyer(i, api)
            task = asyncio.create_task(buyer.lifeCycle(actions))
            tasks.append(task)
        await asyncio.gather(*tasks)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Simulate buyer behaviour')
    parser.add_argument("-n", "--buyers", type=int, default=10, dest='buyers', help='number of buyers')
    parser.add_argument("-a", "--actions" , type=int, default=10, dest='actions', help='number of actions')
    parser.add_argument("-u", "--url" , default="http://localhost:3000", dest='api_url', help='URL of running API')
    args = parser.parse_args()
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main(args.api_url, args.buyers, args.actions))
