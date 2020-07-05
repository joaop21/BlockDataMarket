import aiohttp
import asyncio
import pprint
import random
import api
import sys
import argparse
from datetime import datetime

pp = pprint.PrettyPrinter(indent=2)

class Buyer:

    def __init__(self, id, session):
        self.session = session
        self.category = {}
        self.id = id
        self.queries = []

    async def pickCategory(self):
        categories, _ = await api.getCategories(self.session)
        index = random.randint(0, len(categories)-1)
        self.category = categories[index]
        print("Worker " + str(self.id))

    async def cheaperQuery(self):
        query = random.choice(self.category['possibleQueries'])
        announcements = await api.getAnnouncementsByCategory(self.session, self.category['name'])
        announcementId, price = Buyer.selectCheaperAnnouncement(announcements, query)
        resultQuery = await api.makeQuery(self.session, announcementId, query, price)
        self.queries.append(resultQuery['queryId'])
        printOp("CheaperQuery", self.id, announcementId, query, price)

    async def randomQuery(self):
        announcements, _ = await api.getAllAnnouncements(self.session)
        announcement = random.choice(announcements)
        index = random.randrange(len(announcement['possibleQueries']))
        query = announcement['possibleQueries'][index]
        queryPrice = announcement['prices'][index]
        priceToPay = random.randint(round(queryPrice/2),queryPrice)
        resultQuery = await api.makeQuery(self.session, announcement['announcementId'], query, priceToPay)
        self.queries.append(resultQuery['queryId'])
        printOp("RandomQuery", self.id, announcement['announcementId'], query, priceToPay)

    async def queryLowerThan(self):
        announcements, _ = await api.getAllAnnouncements(self.session)
        query = random.choice(self.category['possibleQueries'])
        meanPrice = await Buyer.getQueryMeanPrice(announcements, query)
        announcements, _ = await api.getAnnouncementsByCategoryLt(self.session, self.category['name'], meanPrice)
        announcement = random.choice(announcements)
        query, price = Buyer.getCheapestQuery(announcement)
        resultQuery = await api.makeQuery(self.session, announcement['announcementId'], query, price)
        self.queries.append(resultQuery['queryId'])
        print("Worker " + str(self.id))
        printOp("QueryLowerThan", self.id, announcement['announcementId'], query, price)

    #alterar id
    async def consultPreviousQueries(self):
        issuedQueries, _ = await api.getQueriesByIssuer(self.session, self.id)
        print("ISSUED " + str(issuedQueries))

    async def consultPreviousQuery(self):
        try:
            queryId = random.choice(self.queries)
            query, _ = await api.getQueryById(self.session, queryId)
        except:
            pass    
        print("Worker " + str(self.id) + " consulted previous query")

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
        pp.pprint(prices)
        return prices[0][0], prices[0][1]
        
    async def lifeCycle(self, numOps):
        pp.pprint("Worker " + str(self.id))
        if not self.category:
            await self.pickCategory()

        for _ in range(numOps):
            try:
                await asyncio.sleep(5)
                op = random.random()
                if op < 0.3:
                    await self.cheaperQuery()
                elif op < 0.6:
                    await self.randomQuery()
                elif op < 0.8:
                    await self.queryLowerThan()
                elif op < 0.9:
                    await self.consultPreviousQueries()
                else:
                    await self.consultPreviousQuery()
            except:
                pass                 

def printOp(op,id, announcementId, query, price):
    now = datetime.now()
    current_time = now.strftime("%H:%M:%S")
    print(current_time + " -> (" + op +") Buyer " + str(id) + " queried " + query + " on announcement " + announcementId + " for " + str(price))

async def main(nBuyers):
    async with aiohttp.ClientSession() as session:
        tasks = []
        for i in range(nBuyers):
            buyer = Buyer(i, session)
            task = asyncio.create_task(buyer.lifeCycle(10))
            tasks.append(task)
        await asyncio.gather(*tasks)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Simulate buyer behaviour')
    parser.add_argument("-n", type=int, default=10, dest='nBuyers', help='number of buyers')
    args = parser.parse_args()
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main(args.nBuyers))
