import os

def create_import_scripts():
    current_dir = os.getcwd()
    
    for folder in os.listdir(current_dir):
        folder_path = os.path.join(current_dir, folder)
        
        if os.path.isdir(folder_path):
            script_path = os.path.join(folder_path, 'import.sh')
            
            content = f'''RESOURCE_ID="12345678-90ab-cdef-1234-56789abcdef0"
terraform import {folder}.this $RESOURCE_ID
'''
            with open(script_path, 'w') as f:
                f.write(content)
            print(f"Created import.sh in {folder}")

# Call the function
create_import_scripts()
