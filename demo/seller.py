import aiohttp
import asyncio
import pprint
import random
import api
import sys
import os
import argparse
import numpy as np
from datetime import datetime

pp = pprint.PrettyPrinter(indent=2)

dir = './docs'


class Seller:

    def __init__(self, id, category, session):
        self.session = session
        self.category = category
        self.id = id
        self.announcement = {}
        self.processedQueries = []

    async def makeAnnouncement(self):
        numQueries = random.randint(1, len(self.category['possibleQueries']))
        queries = random.sample(self.category['possibleQueries'], numQueries)
        filename = random.choice(os.listdir(dir))
        filePath = os.path.join(dir, filename)
        announcement = await api.makeAnnouncement(self.session, filePath, filename, self.category['name'], queries)
        return announcement

    async def updatePrices(self):
        announcementQueries, error = await api.getQueriesByAnnouncement(self.session, self.announcement['announcementId'])
        if not error:
            newQueries = [query for query in announcementQueries if query not in self.processedQueries]
            now = datetime.now()
            current_time = now.strftime("%H:%M:%S")
            print(current_time + " SELLER")
            print("PROCESSED QUERIES " + str(len(self.processedQueries)))
            print("NEW QUERIES " + str(len(newQueries)))
            print("TOTAL " + str(len(announcementQueries)))
            nrQueries = dict(zip(self.announcement['possibleQueries'], np.zeros(len(self.announcement['possibleQueries']))))
            for query in newQueries:
                nrQueries[query['query']]+=1
                self.processedQueries.append(query)
            pp.pprint(nrQueries)
            prices = dict(zip(self.announcement['possibleQueries'], self.announcement['prices']))
            pp.pprint(prices)
            for query in nrQueries.items():
                if(query[1] == 0):
                    prices[query[0]]*=0.9
                else:
                    prices[query[0]]*=(1+ 0.05*query[1])
            pp.pprint(prices)
            self.announcement = await api.updatePrices(self.session, self.announcement['announcementId'], prices)
            pp.pprint(self.announcement)




            



        
    async def lifeCycle(self, numOps):
        pp.pprint("Seller " + str(self.id))
        for _ in range(numOps):
            if not self.announcement:
                announcement = await self.makeAnnouncement()
                if not 'error' in announcement.keys():
                    self.announcement = announcement
            else:
                await asyncio.sleep(10)
                await self.updatePrices()

            '''        
            try:
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
            '''    
async def main(nSellers, categoryName):
    async with aiohttp.ClientSession() as session:
        tasks = []
        category, _ = await api.getCategoryByName(session, categoryName)
        for i in range(nSellers):
            seller = Seller(i, category, session)
            task = asyncio.create_task(seller.lifeCycle(10))
            tasks.append(task)
        await asyncio.gather(*tasks)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Simulate seller behaviour')
    parser.add_argument("-n", type=int, default=10, dest='nSellers', help='number of sellers')
    parser.add_argument("-c", "--category", default="Wikipedia", dest='categoryName', help='name of the category')
    parser.add_argument("-d", "--directory", default="'./docs'", dest='dir', help='path of dir containing data docs')
    args = parser.parse_args()
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main(args.nSellers, args.categoryName))
