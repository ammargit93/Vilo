import os
from flask import Flask, request, jsonify
from werkzeug.utils import secure_filename

app = Flask(__name__)

# print(f"Upload folder: {os.path.abspath()}")  # Debug print
# os.makedirs(UPLOAD_FOLDER, exist_ok=True)

@app.route("/upload", methods=["POST"])
def upload_file():
    try:
        

        if "file" not in request.files:
            return jsonify({"error": "No file part"}), 400

        file = request.files["file"]
        if file.filename == "":
            return jsonify({"error": "No selected file"}), 400

        name = file.filename.split("\\")[-1]
        commithash = file.filename.split("\\")[-2]
        
        print(f"Upload folder: {os.path.abspath(commithash)}")  
        os.makedirs(os.path.abspath(commithash), exist_ok=True)
        
        print("Base name of the file : ",name)
        filename = secure_filename(file.filename)
        save_path = os.path.join(os.path.abspath(commithash), name)
        
        print(f"Saving file to: {save_path}")  
        file.save(save_path)
        
        return jsonify({"message": f"File {file.filename} uploaded successfully!"})
    
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    app.run(debug=True)
