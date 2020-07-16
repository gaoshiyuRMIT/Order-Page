## General Setup
#### Edit configuration files
- Specify these details in `config.json` (refer to `config.example.json`)
    * port number to run the API on
    * database connection details

## API
Make sure you are in `app`: `cd app`

#### Setup
- `go mod download`

#### Usage
- `go run server.go`


## Client App
First make sure you have the API running, and you are in `clientapp`.

#### Setup
- Specify the base URL of your backend API in `config.js` (refer to `config.example.js`)

#### Usage
- All the files in this folder are static, so you can host them on any website hosting platform.
- Open `index.html`.


## Data Import
Make sure you are in `import_script`: `cd import_script`

#### Set up the virtual environment
- create venv: `python3 -m venv venv`
- activate venv: `. venv/bin/activate`
- install libraries: `pip install -r requirements.txt`

#### Usage
- activate venv: - `. venv/bin/activate`
- Make sure the `test_data` folder, `import_data.py` and `config.json` are in the same directory.
- `python import_data.py`