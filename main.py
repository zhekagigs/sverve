import sys

def db_set(key, value):
    with open("database.csv", "a") as f:
       f.write(f"{key},{value}\n")


def db_get(key):
    with open("database.csv", "r") as f:
        value = None
        for line in f:
            line = line.strip()
            if line.startswith(f"{key},"):
                value = line.split(",", 1)[1]
    return value


def main():

    if len(sys.argv) < 2:
        print("Usage: python main.py [get|set] [key] [value]")
        print("  get <key>          - Retrieve value for key")
        print("  set <key> <value>  - Store key-value pair")
        sys.exit(1)

    command = sys.argv[1]

    if command == "get":
        if len(sys.argv) < 3:
            print("Error: get requires a key")
            sys.exit(1)
        key = sys.argv[2]
        result = db_get(key)
        if result:
            print(result)
        else:
            print(f"Key '{key}' not found")

    elif command == "set":
        if len(sys.argv) < 4:
            print("Error: set requires key and value")
            sys.exit(1)
        key = sys.argv[2]
        value = sys.argv[3]
        db_set(key, value)

    else:
        print(f"Unknown command: {command}")
        sys.exit(1)


if __name__ == "__main__":
    main()

