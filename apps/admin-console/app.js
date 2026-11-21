// BRLIA Control - Admin Console JavaScript
// Minimalista, Profissional, Sem Emojis

document.addEventListener('DOMContentLoaded', function() {
    // Initialize the application
    initSidebar();
    initNavigation();
    initUserMenu();
    initPageLoaders();
    initDataFetch();
    initEventListeners();
});

// Sidebar Toggle
function initSidebar() {
    const menuToggle = document.getElementById('menuToggle');
    const sidebar = document.querySelector('.sidebar');
    const mainContent = document.querySelector('.main-content');

    if (menuToggle && sidebar) {
        menuToggle.addEventListener('click', function() {
            sidebar.classList.toggle('open');
            if (sidebar.classList.contains('open')) {
                mainContent.style.marginLeft = '250px';
            } else {
                mainContent.style.marginLeft = '0';
            }
        });

        // Close sidebar when clicking outside
        document.addEventListener('click', function(e) {
            if (!sidebar.contains(e.target) && !menuToggle.contains(e.target)) {
                sidebar.classList.remove('open');
                mainContent.style.marginLeft = '0';
            }
        });
    }
}

// Navigation
function initNavigation() {
    const navItems = document.querySelectorAll('.nav-item');
    const pages = document.querySelectorAll('.page');
    const currentPageLabel = document.getElementById('currentPage');

    navItems.forEach(item => {
        item.addEventListener('click', function(e) {
            e.preventDefault();

            // Remove active class from all nav items
            navItems.forEach(nav => nav.classList.remove('active'));
            
            // Add active class to clicked item
            this.classList.add('active');

            // Hide all pages
            pages.forEach(page => page.classList.remove('active'));

            // Show target page
            const targetPage = this.getAttribute('data-page');
            const target = document.getElementById(`page-${targetPage}`);
            if (target) {
                target.classList.add('active');
                
                // Update current page label
                if (currentPageLabel) {
                    currentPageLabel.textContent = capitalizeFirstLetter(targetPage);
                }

                // Load page data
                loadPageData(targetPage);
            }
        });
    });
}

// User Menu
function initUserMenu() {
    const userButton = document.getElementById('userButton');
    const userDropdown = document.getElementById('userDropdown');

    if (userButton && userDropdown) {
        userButton.addEventListener('click', function(e) {
            e.stopPropagation();
            userDropdown.classList.toggle('show');
        });

        // Close dropdown when clicking outside
        document.addEventListener('click', function(e) {
            if (!userButton.contains(e.target) && !userDropdown.contains(e.target)) {
                userDropdown.classList.remove('show');
            }
        });
    }
}

// Page Loaders
function initPageLoaders() {
    // Initialize any page-specific functionality
}

// Data Fetching
function initDataFetch() {
    // Load initial dashboard data
    loadDashboardData();
    
    // Set up periodic refresh
    setInterval(loadDashboardData, 30000);
}

// Load Dashboard Data
function loadDashboardData() {
    // Simulate API calls
    fetch('/api/v1/organizations?limit=1000').then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Failed to fetch organizations');
    }).then(data => {
        const count = data.total || 0;
        updateStat('orgCount', count);
    }).catch(error => {
        console.error('Error fetching organizations:', error);
    });

    fetch('/api/v1/projects?limit=1000').then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Failed to fetch projects');
    }).then(data => {
        const count = data.total || 0;
        updateStat('projectCount', count);
    }).catch(error => {
        console.error('Error fetching projects:', error);
    });

    fetch('/api/v1/users?limit=1000').then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Failed to fetch users');
    }).then(data => {
        const count = data.total || 0;
        updateStat('userCount', count);
    }).catch(error => {
        console.error('Error fetching users:', error);
    });

    // Load Waterfall data
    loadWaterfallData();
}

// Load Waterfall Data
function loadWaterfallData() {
    fetch('/api/v1/status').then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Failed to fetch Waterfall status');
    }).then(data => {
        if (data.WorkloadsScheduled) {
            updateStat('waterfallCompleted', data.WorkloadsScheduled);
        }
        if (data.WorkloadsCompleted) {
            updateStat('waterfallCompleted', data.WorkloadsCompleted);
        }
        if (data.WorkloadsFailed) {
            updateStat('waterfallFailed', data.WorkloadsFailed);
        }
        if (data.AverageCompletionTime) {
            updateStat('waterfallAvgTime', data.AverageCompletionTime + 'ms');
        }
    }).catch(error => {
        console.error('Error fetching Waterfall status:', error);
    });
}

// Load Page Data
function loadPageData(page) {
    switch(page) {
        case 'organizations':
            loadOrganizations();
            break;
        case 'projects':
            loadProjects();
            break;
        case 'users':
            loadUsers();
            break;
        case 'roles':
            loadRoles();
            break;
        case 'regions':
            loadRegions();
            break;
        case 'waterfall':
            loadWaterfallPage();
            break;
        case 'workloads':
            loadWorkloads();
            break;
    }
}

// Load Organizations
function loadOrganizations() {
    fetch('/api/v1/organizations').then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Failed to fetch organizations');
    }).then(data => {
        const tbody = document.getElementById('orgTableBody');
        if (tbody) {
            tbody.innerHTML = '';
            data.organizations.forEach(org => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${org.id.substring(0, 8)}</td>
                    <td>${org.name}</td>
                    <td>${org.description || '-'}</td>
                    <td><span class="status-badge ${org.status}">${capitalizeFirstLetter(org.status)}</span></td>
                    <td>${org.project_count || 0}</td>
                    <td>${org.user_count || 0}</td>
                    <td>
                        <button class="action-btn view" title="Visualizar"></button>
                        <button class="action-btn edit" title="Editar"></button>
                        <button class="action-btn delete" title="Excluir"></button>
                    </td>
                `;
                tbody.appendChild(row);
            });
        }
    }).catch(error => {
        console.error('Error fetching organizations:', error);
    });
}

// Load Projects
function loadProjects() {
    fetch('/api/v1/projects').then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Failed to fetch projects');
    }).then(data => {
        const tbody = document.getElementById('projectTableBody');
        if (tbody) {
            tbody.innerHTML = '';
            data.projects.forEach(project => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${project.id.substring(0, 8)}</td>
                    <td>${project.name}</td>
                    <td>${project.organization_name || '-'}</td>
                    <td>${project.region_name || '-'}</td>
                    <td><span class="status-badge ${project.status}">${capitalizeFirstLetter(project.status)}</span></td>
                    <td>${project.resource_count || 0}</td>
                    <td>
                        <button class="action-btn view" title="Visualizar"></button>
                        <button class="action-btn edit" title="Editar"></button>
                        <button class="action-btn delete" title="Excluir"></button>
                    </td>
                `;
                tbody.appendChild(row);
            });
        }
    }).catch(error => {
        console.error('Error fetching projects:', error);
    });
}

// Load Users
function loadUsers() {
    fetch('/api/v1/users').then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Failed to fetch users');
    }).then(data => {
        const tbody = document.getElementById('userTableBody');
        if (tbody) {
            tbody.innerHTML = '';
            data.users.forEach(user => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${user.id.substring(0, 8)}</td>
                    <td>${user.display_name}</td>
                    <td>${user.username}</td>
                    <td>${user.email}</td>
                    <td>${user.organization_name || '-'}</td>
                    <td><span class="status-badge ${user.status}">${capitalizeFirstLetter(user.status)}</span></td>
                    <td>
                        <button class="action-btn view" title="Visualizar"></button>
                        <button class="action-btn edit" title="Editar"></button>
                        <button class="action-btn delete" title="Excluir"></button>
                    </td>
                `;
                tbody.appendChild(row);
            });
        }
    }).catch(error => {
        console.error('Error fetching users:', error);
    });
}

// Load Roles
function loadRoles() {
    fetch('/api/v1/roles').then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Failed to fetch roles');
    }).then(data => {
        const tbody = document.getElementById('roleTableBody');
        if (tbody) {
            tbody.innerHTML = '';
            data.roles.forEach(role => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${role.id.substring(0, 8)}</td>
                    <td>${role.name}</td>
                    <td>${role.description || '-'}</td>
                    <td>${role.permissions.length}</td>
                    <td><span class="status-badge ${role.is_system ? 'active' : 'inactive'}">${role.is_system ? 'Sim' : 'Não'}</span></td>
                    <td>
                        <button class="action-btn view" title="Visualizar"></button>
                        <button class="action-btn edit" title="Editar"></button>
                        ${!role.is_system ? '<button class="action-btn delete" title="Excluir"></button>' : ''}
                    </td>
                `;
                tbody.appendChild(row);
            });
        }
    }).catch(error => {
        console.error('Error fetching roles:', error);
    });
}

// Load Regions
function loadRegions() {
    fetch('/api/v1/regions').then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Failed to fetch regions');
    }).then(data => {
        const tbody = document.getElementById('regionTableBody');
        if (tbody) {
            tbody.innerHTML = '';
            data.regions.forEach(region => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${region.id.substring(0, 8)}</td>
                    <td>${region.name}</td>
                    <td>${region.location}</td>
                    <td>${region.country}</td>
                    <td><span class="status-badge ${region.status}">${capitalizeFirstLetter(region.status)}</span></td>
                    <td>${region.zone_count || 0}</td>
                    <td>${region.cluster_count || 0}</td>
                    <td>
                        <button class="action-btn view" title="Visualizar"></button>
                        <button class="action-btn edit" title="Editar"></button>
                    </td>
                `;
                tbody.appendChild(row);
            });
        }
    }).catch(error => {
        console.error('Error fetching regions:', error);
    });
}

// Load Waterfall Page
function loadWaterfallPage() {
    // Load Waterfall engine status
    fetch('/api/v1/status').then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Failed to fetch Waterfall status');
    }).then(data => {
        updateStat('waterfallActive', data.ActiveWorkloads || 0);
        updateStat('waterfallPending', data.PendingWorkloads || 0);
        updateStat('waterfallCompleted', data.WorkloadsCompleted || 0);
        updateStat('waterfallAvgTime', (data.AverageCompletionTime || 0) + 'ms');
        
        // Update engine status
        const engineStatusEl = document.getElementById('engineStatus');
        if (engineStatusEl) {
            engineStatusEl.textContent = data.Status || 'Operacional';
        }
        
        const engineWorkersEl = document.getElementById('engineWorkers');
        if (engineWorkersEl) {
            engineWorkersEl.textContent = data.WorkerCount || 10;
        }
        
        const engineStartedEl = document.getElementById('engineStarted');
        if (engineStartedEl) {
            engineStartedEl.textContent = formatDate(data.StartedAt) || '2026-11-21 10:00:00';
        }
    }).catch(error => {
        console.error('Error fetching Waterfall status:', error);
    });
    
    // Load recent workloads
    fetch('/api/v1/workloads?limit=10').then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Failed to fetch workloads');
    }).then(data => {
        const tbody = document.getElementById('waterfallWorkloads');
        if (tbody) {
            tbody.innerHTML = '';
            data.workloads.forEach(workload => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${workload.id.substring(0, 8)}</td>
                    <td>${workload.name}</td>
                    <td>${workload.type}</td>
                    <td><span class="status-badge ${workload.status}">${capitalizeFirstLetter(workload.status)}</span></td>
                    <td>${formatPriority(workload.priority)}</td>
                    <td>${workload.cluster_name || '-'}</td>
                    <td>${formatDuration(workload.execution_time)}</td>
                `;
                tbody.appendChild(row);
            });
        }
    }).catch(error => {
        console.error('Error fetching workloads:', error);
    });
}

// Load Workloads
function loadWorkloads() {
    fetch('/api/v1/workloads').then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Failed to fetch workloads');
    }).then(data => {
        // Handle workloads data
    }).catch(error => {
        console.error('Error fetching workloads:', error);
    });
}

// Event Listeners
function initEventListeners() {
    // Create organization button
    const createOrgBtn = document.getElementById('createOrgBtn');
    if (createOrgBtn) {
        createOrgBtn.addEventListener('click', function() {
            showCreateOrganizationModal();
        });
    }

    // Create project button
    const createProjectBtn = document.getElementById('createProjectBtn');
    if (createProjectBtn) {
        createProjectBtn.addEventListener('click', function() {
            showCreateProjectModal();
        });
    }

    // Create user button
    const createUserBtn = document.getElementById('createUserBtn');
    if (createUserBtn) {
        createUserBtn.addEventListener('click', function() {
            showCreateUserModal();
        });
    }

    // Create role button
    const createRoleBtn = document.getElementById('createRoleBtn');
    if (createRoleBtn) {
        createRoleBtn.addEventListener('click', function() {
            showCreateRoleModal();
        });
    }

    // Create region button
    const createRegionBtn = document.getElementById('createRegionBtn');
    if (createRegionBtn) {
        createRegionBtn.addEventListener('click', function() {
            showCreateRegionModal();
        });
    }

    // Search functionality
    const orgSearch = document.getElementById('orgSearch');
    if (orgSearch) {
        orgSearch.addEventListener('input', function() {
            debounceSearch('organizations', this.value);
        });
    }

    const projectSearch = document.getElementById('projectSearch');
    if (projectSearch) {
        projectSearch.addEventListener('input', function() {
            debounceSearch('projects', this.value);
        });
    }

    const userSearch = document.getElementById('userSearch');
    if (userSearch) {
        userSearch.addEventListener('input', function() {
            debounceSearch('users', this.value);
        });
    }

    const roleSearch = document.getElementById('roleSearch');
    if (roleSearch) {
        roleSearch.addEventListener('input', function() {
            debounceSearch('roles', this.value);
        });
    }
}

// Search debounce
let searchTimeout;
function debounceSearch(page, query) {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(function() {
        search(page, query);
    }, 500);
}

// Search
function search(page, query) {
    // Implement search functionality
    console.log('Searching', page, 'for', query);
}

// Helper Functions
function updateStat(elementId, value) {
    const element = document.getElementById(elementId);
    if (element) {
        element.textContent = value;
    }
}

function capitalizeFirstLetter(string) {
    if (!string) return '';
    return string.charAt(0).toUpperCase() + string.slice(1);
}

function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleString('pt-BR');
}

function formatDuration(durationMs) {
    if (!durationMs) return '-';
    const seconds = Math.floor(durationMs / 1000);
    if (seconds < 60) {
        return seconds + 's';
    }
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    if (minutes < 60) {
        return minutes + 'm ' + remainingSeconds + 's';
    }
    const hours = Math.floor(minutes / 60);
    const remainingMinutes = minutes % 60;
    return hours + 'h ' + remainingMinutes + 'm';
}

function formatPriority(priority) {
    const priorities = ['Lowest', 'Low', 'Normal', 'High', 'Highest'];
    return priorities[priority] || 'Normal';
}

// Modal Functions
function showCreateOrganizationModal() {
    alert('Criar Nova Organização - Implementar modal');
}

function showCreateProjectModal() {
    alert('Criar Novo Projeto - Implementar modal');
}

function showCreateUserModal() {
    alert('Criar Novo Usuário - Implementar modal');
}

function showCreateRoleModal() {
    alert('Criar Nova Role - Implementar modal');
}

function showCreateRegionModal() {
    alert('Criar Nova Região - Implementar modal');
}
