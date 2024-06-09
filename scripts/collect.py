import argparse

parser = argparse.ArgumentParser(description="Collect the meaningful data on the web")
parser.add_argument("-k","--keyword", type=str, help="set the keywords to research in deep, it would be your brand or product names separated by a comma")
parser.add_argument("-n", "--number", type=int, help="set the number of blogs to retrieve", default=3) 
parser.add_argument("-p", "--platform", type=str, help="social-media or search-engine to search keywords from (choose one of from: naverblog)",choices=['naverblog', 'wip:navercafe'])
parser.add_argument("-q", "--quiet", action='store_true',help="set the quiet mode to suppress debug message from command line output")

args = parser.parse_args()
print(args.quiet)
