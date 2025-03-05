document.addEventListener("DOMContentLoaded", function() {
  // =========================
  // Country Info Lookup
  // =========================
  var infoBtn = document.getElementById("getCountryBtn");
  infoBtn.addEventListener("click", function() {
    var code = document.getElementById("countryCode").value.trim().toLowerCase();
    var cityLimit = document.getElementById("cityLimit").value.trim();
    
    // Validate that a proper 2-letter country code is provided.
    if (!code || code.length !== 2) {
      alert("Please enter a valid 2-letter country code.");
      return;
    }
    
    // Build the API URL for the info endpoint.
    var url = "/countryinfo/v1/info/" + code;
    if (cityLimit !== "") {
      url += "?limit=" + encodeURIComponent(cityLimit);
    }
    
    // Fetch country info from the server.
    fetch(url)
      .then(function(response) {
        if (!response.ok) {
          throw new Error("Country info fetch failed: " + response.status);
        }
        return response.json();
      })
      .then(function(data) {
        var infoResultDiv = document.getElementById("infoResult");
        infoResultDiv.innerHTML =
          "<h3>Country Info for " + code.toUpperCase() + "</h3>" +
          "<pre>" + JSON.stringify(data, null, 2) + "</pre>";
      })
      .catch(function(err) {
        console.error("Error fetching country info:", err);
        alert("Error fetching country info. See console for details.");
      });
  });

  // =========================
  // Population Data Lookup
  // =========================
  var popBtn = document.getElementById("getPopulationBtn");
  popBtn.addEventListener("click", function() {
    var code = document.getElementById("popCountryCode").value.trim().toLowerCase();
    var startYear = document.getElementById("startYear").value.trim();
    var endYear = document.getElementById("endYear").value.trim();
    
    // Validate that a proper 2-letter country code is provided.
    if (!code || code.length !== 2) {
      alert("Please enter a valid 2-letter country code.");
      return;
    }
    
    // Build the API URL for the population endpoint.
    var url = "/countryinfo/v1/population/" + code;
    if (startYear !== "" && endYear !== "") {
      url += "?limit=" + encodeURIComponent(startYear + "-" + endYear);
    }
    
    // Fetch population data from the server.
    fetch(url)
      .then(function(response) {
        if (!response.ok) {
          throw new Error("Population fetch failed: " + response.status);
        }
        return response.json();
      })
      .then(function(data) {
        var popResultDiv = document.getElementById("popResult");
        popResultDiv.innerHTML =
          "<h3>Population Data for " + code.toUpperCase() + "</h3>" +
          "<pre>" + JSON.stringify(data, null, 2) + "</pre>";
      })
      .catch(function(err) {
        console.error("Error fetching population data:", err);
        alert("Error fetching population data. See console for details.");
      });
  });

  // =========================
  // Diagnostics Lookup
  // =========================
  var diagBtn = document.getElementById("getDiagnosticsBtn");
  diagBtn.addEventListener("click", function() {
    // Fetch diagnostics data from the server.
    fetch("/countryinfo/v1/status/")
      .then(function(response) {
        if (!response.ok) {
          throw new Error("Diagnostics fetch failed: " + response.status);
        }
        return response.json();
      })
      .then(function(data) {
        var diagResultDiv = document.getElementById("diagResult");
        diagResultDiv.innerHTML =
          "<h3>Diagnostics Status</h3>" +
          "<pre>" + JSON.stringify(data, null, 2) + "</pre>";
      })
      .catch(function(err) {
        console.error("Error fetching diagnostics data:", err);
        alert("Error fetching diagnostics data. See console for details.");
      });
  });
});
