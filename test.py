from flask import Flask, jsonify, request
import os
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.backends import default_backend
import zlib  
import io

app = Flask(__name__)

UPLOAD_FOLDER = "uploads"
os.makedirs(UPLOAD_FOLDER, exist_ok=True)

KEY = b"thisis32bytekeythisis32bytekey!!"

def decrypt_and_decompress(encrypted_file_data, key):
    backend = default_backend()
    iv = encrypted_file_data[:16] 
    encrypted_data = encrypted_file_data[16:] 
    
    cipher = Cipher(algorithms.AES(key), modes.CFB(iv), backend=backend)
    decryptor = cipher.decryptor()
    decrypted_data = decryptor.update(encrypted_data) + decryptor.finalize()
    decompressed_data = zlib.decompress(decrypted_data)
    return decompressed_data

@app.route("/hello", methods=["GET", "POST"])
def hello_world():
    if request.method == "POST":
        # Check if any files are uploaded
        if "file" not in request.files:
            return jsonify({"error": "No file part"}), 400
        
        files = request.files.getlist("file")  # Get all files from the request

        saved_files = []
        for file in files:
            if file.filename == "":
                return jsonify({"error": "No selected file"}), 400
            
            # Read file content
            encrypted_content = file.read()

            try:
                # Decrypt and decompress
                processed_content = decrypt_and_decompress(encrypted_content, KEY)
                
                # Save the processed content to the upload folder
                file_path = os.path.join(UPLOAD_FOLDER, file.filename.replace(".enc", ""))
                with open(file_path, "wb") as f:
                    f.write(processed_content)
                
                saved_files.append(file.filename.replace(".enc", ""))
            except Exception as e:
                return jsonify({"error": f"Failed to process {file.filename}: {str(e)}"}), 500

        return jsonify({"message": "Files uploaded, decrypted, and decompressed successfully", "files": saved_files})
    
    return jsonify({"message": "Hello"})

if __name__ == "__main__":
    app.run(debug=True)
