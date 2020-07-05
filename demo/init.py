import aiohttp
import argparse
import asyncio
import pprint
from api import API

pp = pprint.PrettyPrinter(indent=2)


async def main(api_url, name):
    async with aiohttp.ClientSession() as session:
        api = API(api_url, session)
        id = await api.makeIdentification(args.name)
        pp.pprint(id)
        category = await api.makeCategory("Wikipedia", "['Title', 'Abstract', 'Url',  'Sections']")
        pp.pprint(category)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Create an identification and category')
    parser.add_argument("-n", "--name", default="Default Org", dest='name', help='identification name')
    parser.add_argument("-u", "--url" , default="http://localhost:3000", dest='api_url', help='URL of running API')
    args = parser.parse_args()
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main(args.api_url, args.name))


