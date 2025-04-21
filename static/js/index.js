document.addEventListener("DOMContentLoaded", function () {
    const modals = {
      add: document.getElementById("addModal"),
      delete: document.getElementById("deleteModal"),
      search: document.getElementById("searchModal")
    };
  
    const buttons = {
      add: document.getElementById("addBtn"),
      delete: document.getElementById("deleteBtn"),
      search: document.getElementById("searchBtn")
    };
  
    const closes = {
      add: document.getElementById("addClose"),
      delete: document.getElementById("deleteClose"),
      search: document.getElementById("searchClose")
    };
  
    if (buttons.add && closes.add) {
      buttons.add.onclick = () => modals.add.style.display = "block";
      closes.add.onclick = () => modals.add.style.display = "none";
    }
  
    if (buttons.delete && closes.delete) {
      buttons.delete.onclick = () => modals.delete.style.display = "block";
      closes.delete.onclick = () => modals.delete.style.display = "none";
    }
  
    if (buttons.search && closes.search) {
      buttons.search.onclick = () => modals.search.style.display = "block";
      closes.search.onclick = () => modals.search.style.display = "none";
    }
  
    window.onclick = function(event) {
      for (let key in modals) {
        if (event.target === modals[key]) {
          modals[key].style.display = "none";
        }
      }
    };
  
    const orgInput = document.getElementById("orgInput");
    const cityInput = document.getElementById("cityInput");
    const phoneInput = document.getElementById("phoneInput");
    const saveBtn = document.getElementById("saveBtn");
  
    function validateForm() {
      const phoneRegex = /^\+7-\d{3}-\d{3}-\d{2}-\d{2}$/;
      const isValidPhone = phoneRegex.test(phoneInput.value);
      phoneInput.classList.toggle('invalid', !isValidPhone);
      const isFilled = orgInput.value.trim() !== "" && cityInput.value.trim() !== "";
      saveBtn.disabled = !(isFilled && isValidPhone);
    }
  
    [orgInput, cityInput, phoneInput].forEach(input => {
      input.addEventListener("input", validateForm);
    });
  });