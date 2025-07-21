export interface Blog {
    id: string
    title: string
    summary: string

    createdAt: string
    updatedAt: string

    prequelId: string | null // Blog ID of prequel (nullable)
    sequelId: string | null// Blog ID of sequel (nullable)

    parts: Array<string>// List of Blog IDs if this is part of a multipart series
}
