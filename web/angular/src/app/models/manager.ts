export interface Manager {
    id: number;
    username: string
    tracker_admin: boolean
    register_time: number
}

export class ManagerRoleOnProject {
    manager: Manager;
    role: number;

    public static fromEntity(data: { role: number, manager: Manager }): ManagerRoleOnProject {
        let m = new ManagerRoleOnProject();
        m.role = data.role;
        m.manager = data.manager;
        return m;
    }

    get readRole(): boolean {
        return (this.role & 1) != 0
    }

    set readRole(role: boolean) {
        if (role) {
            this.role |= 1
        } else {
            this.role &= ~1
        }
    }

    get editRole(): boolean {
        return (this.role & 2) != 0
    }

    set editRole(role: boolean) {
        if (role) {
            this.role |= 2
        } else {
            this.role &= ~2
        }
    }

    get manageRole(): boolean {
        return (this.role & 4) != 0
    }

    set manageRole(role: boolean) {
        if (role) {
            this.role |= 4
        } else {
            this.role &= ~4
        }
    }
}
