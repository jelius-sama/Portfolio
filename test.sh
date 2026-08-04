#!/bin/bash

# Base URL
BASE_URL="http://localhost:6969"
MD_FILE="$HOME/Documents/1a42b1d.md"

# Array of blog posts
declare -a BLOGS=(
    "caveats-workarounds-openmediacloud|Caveats & Workarounds — OpenMediaCloud|Architectural Workarounds for Compute-Dependent Features"
    "setup-guide-openmediacloud|Setup Guide — OpenMediaCloud|This guide walks you through setting up OpenMediaCloud on a Linux server"
    "proxy-cuts-jellyfin-hosting-costs|A Proxy that cuts my Jellyfin Hosting costs by 80%|Say goodbye to egress costs"
    "safari-is-a-dumpster-fire|Safari Is a Dumpster Fire and Apple Knows It|Safari is not \"privacy-focused brilliance.\" — It's a broken, inconsistent, developer-hostile piece of shit."
    "convert-cbz|Convert Folders to CBZ Archives with \"convert-cbz\" — Fast, Cross-Platform Comic Converter|A powerful, cross-platform tool for converting directory of folders of images into CBZ comic book archives. \"convert-cbz\" helps you to quickly create clean, organized CBZ files for your entire comic or manga library."
    "email-done-right-part-3|Email Done Right: My End-to-End Journey from Gmail to AWS SES — Part 3|Getting SES production access, cleaning up configs, and finally sending emails that look truly professional."
    "email-done-right-part-2|Email Done Right: My End-to-End Journey from Gmail to AWS SES — Part 2|Setting up AWS SES, DNS records, and testing email flow while dealing with sandbox restrictions."
    "email-done-right-part-1|Email Done Right: My End-to-End Journey from Gmail to AWS SES — Part 1|Why Gmail's \"sent via\" headers and limits pushed me to find a professional email-sending solution."
    "hello-world-im-jelius|Hello World — I'm Jelius|A quick introduction to who I am, what I do, and why I build. Meet the mind behind jelius.dev."
)

# Post each blog
for BLOG in "${BLOGS[@]}"; do
    IFS='|' read -r ID TITLE EXCERPT <<< "$BLOG"
    
    echo "Posting: $TITLE"
    
    curl -X POST "$BASE_URL/api/blog" \
        -F "markdown=@$MD_FILE" \
        -F "title=$TITLE" \
        -F "excerpt=$EXCERPT" \
        -F "prequel_id=" \
        -F "sequel_id="
    
    echo -e "\n---\n"
done
