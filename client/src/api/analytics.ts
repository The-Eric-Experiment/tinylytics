export interface Summaries {
    pageViews: number;
    sessions: number;
}

export function getSummaries(domain: string): Promise<Summaries> {
    return fetch(`/api/${domain}/summaries?p=alltime`).then(o => o.json());
}