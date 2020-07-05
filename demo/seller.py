import aiohttp
import asyncio
import pprint
import random
from api import API
import sys
import os
import argparse
import numpy as np
from datetime import datetime

pp = pprint.PrettyPrinter(indent=2)

dir = './docs'


class Seller:

    def __init__(self, id, category, api):
        self.api = api
        self.category = category
        self.id = id
        self.announcement = {}
        self.processedQueries = []

    async def makeAnnouncement(self):
        numQueries = random.randint(1, len(self.category['possibleQueries']))
        queries = random.sample(self.category['possibleQueries'], numQueries)
        filename = random.choice(os.listdir(dir))
        filePath = os.path.join(dir, filename)
        announcement = await self.api.makeAnnouncement(filePath, filename, self.category['name'], queries)
        return announcement

    async def updatePrices(self):
        announcementQueries, error = await self.api.getQueriesByAnnouncement(self.announcement['announcementId'])
        if not error:
            newQueries = [query for query in announcementQueries if query not in self.processedQueries]
            nrQueries = dict(zip(self.announcement['possibleQueries'], np.zeros(len(self.announcement['possibleQueries']))))
            for query in newQueries:
                nrQueries[query['query']]+=1
                self.processedQueries.append(query)
            prices = dict(zip(self.announcement['possibleQueries'], self.announcement['prices']))
            old_prices = prices.copy()
            for query in nrQueries.items():
                if(query[1] == 0):
                    prices[query[0]]*=0.9
                else:
                    prices[query[0]]*=(1+ 0.05*query[1])
            self.announcement = await self.api.updatePrices(self.announcement['announcementId'], prices)
            printOp(self.id, self.announcement['announcementId'], old_prices, prices)

    async def lifeCycle(self, numOps):
        for _ in range(numOps):
            if not self.announcement:
                announcement = await self.makeAnnouncement()
                if not 'error' in announcement.keys():
                    self.announcement = announcement
            else:
                await asyncio.sleep(10)
                await self.updatePrices() 
            

def printOp(id, announcementId, old_prices, new_prices):
    now = datetime.now()
    current_time = now.strftime("%H:%M:%S")
    print(current_time + " -> Seller " + str(id) + " changed prices on announcement " + announcementId)
    pp.pprint("Old Prices: " + str(old_prices))
    pp.pprint("New Prices: " + str(new_prices))
        
async def main(api_url, sellers, actions, categoryName):
    async with aiohttp.ClientSession() as session:
        api = API(api_url, session)
        tasks = []
        category, _ = await api.getCategoryByName(categoryName)
        for i in range(sellers):
            seller = Seller(i, category, api)
            task = asyncio.create_task(seller.lifeCycle(actions))
            tasks.append(task)
        await asyncio.gather(*tasks)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Simulate seller behaviour')
    parser.add_argument("-n", "-sellers", type=int, default=10, dest='sellers', help='number of sellers')
    parser.add_argument("-c", "--category", default="Wikipedia", dest='categoryName', help='name of the category')
    parser.add_argument("-d", "--directory", default="./docs", dest='dir', help='path of dir containing data docs')
    parser.add_argument("-a", "--actions", type=int, default=10, dest='actions', help='number of actions')
    parser.add_argument("-u", "--url" , default="http://localhost:3000", dest='api_url', help='URL of running API')
    args = parser.parse_args()
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main(args.api_url, args.sellers, args.actions, args.categoryName))
