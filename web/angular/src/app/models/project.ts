export interface Project {
    id: number;
    priority: number;
    motd: string;
    name: string;
    clone_url: string;
    git_repo: string;
    version: string;
    public: boolean;
    chain: number;
    hidden: boolean;
    paused: boolean;
    assign_rate: number;
    submit_rate: number;
}
