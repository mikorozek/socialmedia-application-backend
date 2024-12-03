import requests

def main():
    print("Hello, World!")
    response = requests.get('https://jsonplaceholder.typicode.com/todos/1')
    print(response.json())
    print("cos")
if __name__ == '__main__':
    main()
