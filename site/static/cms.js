const cms_root_fields = JSON.parse(document.getElementById('p-root-data').textContent);

// Parse p-scope attribute from the html element to get local fields
function parsePScope() {
    const htmlElement = document.documentElement;
    const pScope = htmlElement.getAttribute('p-scope');
    if (!pScope) return {};
    
    const fields = {};
    // Split by semicolon to get individual assignments
    const assignments = pScope.split(';');
    
    for (const assignment of assignments) {
        const trimmed = assignment.trim();
        if (!trimmed) continue;
        
        // Extract variable name (everything before the first =)
        const eqIndex = trimmed.indexOf('=');
        if (eqIndex === -1) continue;
        
        const varName = trimmed.substring(0, eqIndex).trim();
        if (varName) {
            // Initialize with empty string - actual values come from DOM bindings
            fields[varName] = '';
        }
    }
    
    return fields;
}

const cms_local_fields = parsePScope();

function createInputs(obj, container, title) {
    container.appendChild(document.createElement('br'));
    fieldset = document.createElement('fieldset');
    legend = document.createElement('legend');
    legend.textContent=title;
    fieldset.appendChild(legend);
    for (const key in obj) {
        if (Object.hasOwnProperty.call(obj, key)) {
            const label = document.createElement('label');
            label.htmlFor = key;
            label.textContent = key;
            const div = document.createElement('div').appendChild(label).parentNode;
            fieldset.appendChild(div);

            if (Array.isArray(obj[key])) {
                // Create a hidden input with p-model to bind to Pattr
                const hiddenInput = document.createElement('input');
                hiddenInput.type = 'hidden';
                hiddenInput.id = key;
                hiddenInput.name = key;
                hiddenInput.setAttribute('p-model', key);
                hiddenInput.value = obj[key].join(',');
                fieldset.appendChild(hiddenInput);
                
                // Create a container for array items
                const arrayContainer = document.createElement('div');
                arrayContainer.id = key + '-array-container';
                arrayContainer.style.marginLeft = '10px';
                
                // Create input for each array item
                obj[key].forEach((item, index) => {
                    const itemDiv = document.createElement('div');
                    itemDiv.style.display = 'flex';
                    itemDiv.style.alignItems = 'center';
                    itemDiv.style.marginBottom = '5px';
                    
                    const input = document.createElement('input');
                    input.type = 'text';
                    input.value = item;
                    input.dataset.arrayKey = key;
                    input.dataset.arrayIndex = index;
                    input.style.marginRight = '5px';
                    input.addEventListener('input', function() {
                        updateArrayFromInputs(key);
                    });
                    
                    const removeBtn = document.createElement('button');
                    removeBtn.type = 'button';
                    removeBtn.textContent = '×';
                    removeBtn.style.marginLeft = '5px';
                    removeBtn.addEventListener('click', function() {
                        itemDiv.remove();
                        updateArrayFromInputs(key);
                    });
                    
                    itemDiv.appendChild(input);
                    itemDiv.appendChild(removeBtn);
                    arrayContainer.appendChild(itemDiv);
                });
                
                // Add button to add new items
                const addBtn = document.createElement('button');
                addBtn.type = 'button';
                addBtn.textContent = '+ Add item';
                addBtn.addEventListener('click', function() {
                    const index = arrayContainer.querySelectorAll('input[data-array-key="' + key + '"]').length;
                    const itemDiv = document.createElement('div');
                    itemDiv.style.display = 'flex';
                    itemDiv.style.alignItems = 'center';
                    itemDiv.style.marginBottom = '5px';
                    
                    const input = document.createElement('input');
                    input.type = 'text';
                    input.value = '';
                    input.dataset.arrayKey = key;
                    input.dataset.arrayIndex = index;
                    input.style.marginRight = '5px';
                    input.addEventListener('input', function() {
                        updateArrayFromInputs(key);
                    });
                    
                    const removeBtn = document.createElement('button');
                    removeBtn.type = 'button';
                    removeBtn.textContent = '×';
                    removeBtn.style.marginLeft = '5px';
                    removeBtn.addEventListener('click', function() {
                        itemDiv.remove();
                        updateArrayFromInputs(key);
                    });
                    
                    itemDiv.appendChild(input);
                    itemDiv.appendChild(removeBtn);
                    arrayContainer.insertBefore(itemDiv, addBtn);
                    input.focus();
                });
                
                arrayContainer.appendChild(addBtn);
                fieldset.appendChild(arrayContainer);
            } else {
                const input = document.createElement('input');
                input.type = 'text';
                let attribute = "p-model";
                if (typeof obj[key] === 'number') {
                    input.type = 'number';
                }
                input.id = key;
                input.name = key;
                input.placeholder = key;
                input.setAttribute(attribute, key);
                fieldset.appendChild(input);
            }
            
            fieldset.appendChild(document.createElement('br'));
            fieldset.appendChild(document.createElement('br'));
        }
    }
    container.appendChild(fieldset);
}

// Helper function to update the array from individual inputs
function updateArrayFromInputs(key) {
    const inputs = document.querySelectorAll('input[data-array-key="' + key + '"]');
    const newArray = Array.from(inputs).map(input => input.value).filter(v => v !== '');
    
    // Update the hidden input and dispatch input event for Pattr
    const hiddenInput = document.getElementById(key);
    if (hiddenInput) {
        hiddenInput.value = newArray.join(',');
        hiddenInput.dispatchEvent(new Event('input', { bubbles: true }));
    }
}
const cms = document.getElementById('plenti_cms');
createInputs(cms_root_fields, cms, "Root Data");
createInputs(cms_local_fields, cms, "Local Data");

document.getElementById('toggle_plenti_cms').addEventListener('click', function () {
    const menu = document.getElementById('plenti_cms');
    menu.classList.toggle('menu-visible');
});
