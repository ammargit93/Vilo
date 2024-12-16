import os
from flask import Flask, request, jsonify

app = Flask(__name__)

UPLOAD_FOLDER = os.path.abspath("uploaded_files")  # Use absolute path
print(f"Upload folder: {UPLOAD_FOLDER}")  # Debug print
os.makedirs(UPLOAD_FOLDER, exist_ok=True)

@app.route("/upload", methods=["POST"])
def upload_file():
    try:
        print(f"Current working directory: {os.getcwd()}")  # Debug working directory

        if "file" not in request.files:
            return jsonify({"error": "No file part"}), 400
        
        file = request.files["file"]
        if file.filename == "":
            return jsonify({"error": "No selected file"}), 400
        
        save_path = os.path.join(UPLOAD_FOLDER, file.filename)
        print(f"Saving file to: {save_path}")  # Debug save path
        file.save(save_path)
        
        return jsonify({"message": f"File {file.filename} uploaded successfully!"})
    
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    app.run(debug=True)
