# AI Usage Report

> Honestly, I could not charge my Claude Code, so I had to develop big parts of this project with my personal knowledge. 
> I used AI in most important sections like learning another strategy in Redis implementation and comparing it with my own solution.

---

## 1. Tools & Models

List every AI tool/model you used, even if only briefly.

| Tool / Product | Model | Purpose | Frequency |
|---|---|---|---|
| Claude | Sonnet 5 | Researching, Documentation | Occasional |
| ChatGPT | GPT-5.6 Luna | Researching | A few times |
| Cursor | Composer 2.5 | Implementation, Testing, Code Review | Continuous |

---

## 2. Stages of Development

### Planning / Requirements Analysis
- No, I did not use AI for designing architecture or breaking tasks. The system is simple and no need to AI. In this stage I used AI for research about strategies for storing data in Redis. I know and have worked with Sorted Set (ZSet) and I researched about If there was another option or not, so it represents HyperLogLog but for simplicity in implementation at this time I chose Sorted Set.

- I added the description of task + my 2 type design (I described with ASCII diagram text based - It uses the less token) in Claude project context for feature purposes like documentation.

### Scaffolding / Boilerplate
- I worked with most of popular architectures like Clean Architecture, Hexagonal , DDD, so according to the conventions I created structures of my previous projects. I used my favorite and standard structure of Clean Architecture.

- I used AI just for rewriting (removing unused sections from copied) some handy files like Makefile, docker-compose.yml, Dockerfile.

### Implementation
- I always care about the domain layer in all of my projects, so I personally write the domain layer to ensure correctness. I write the types in domain, factory with validation is important for me and I write it.

- The storage layer is important layer in this case, so I define the interfaces and write critical sections and use autocompletion for library methods.

- I used Cursor + GitHub integration for improvement if needed. But for fake storage that just needed for testing purposes, I used AI to generate that methods according to the real storage and interface. 

### Debugging
- No, I did not use AI because there was not unexped bug. I developed this project according Clean Architecture layers and at each step I test the implemented part. I mean I developed the project layer by layer and testing the layer immediately.

### Testing
- Yes, writing test is one of most sections that I use writing tests every time.  The scenarios are important and implementation is not complex thing and most of the time is simple code.

- I define the scenarios for AI and give some tips and then I review all implemented tests and check the scenarios or wrong data in test tables, or add some data.

### Documentation
- Yes, I always use cheap and free AI models for generating documents, because it saves time.

- Especially  I used in generating .md files. I give my concepts and describe some parts and the sorting of data, after generating, I read all data and I personally edit it if there is wrong data.

- One of good parts is generating mermaid in .md files. It takes time but AI does it in fastest way.

### Code Review / Refactoring
- I always break down my functions to reusable parts if needed and do it incrementally. So, sometimes ask AI do it and save my time. 

- Time complexity is important in large scales, so sometimes I asked AI for checing the important parts of the system for performance and improvements.

---

## 3. Prompts

### Prompt #1
- **Stage:** Researching
- **Prompt:**
  ```
  I need to store (user_id, segment) pairs in Redis
  1. Expire after 14 days
  2. Scale: millions of users and hundreds of segments

  Is redis sorted set good fit or is there another approach for this scale?
  ```
- **Result:** It introduce HyperLogLog as alternative approach for very large segments where approximate count is acceptable

### Prompt #2
- **Stage:** Researching
- **Prompt:**
  ```
  Compare Sorted Set vs HyperLogLog for this use case and give me the trade-offs
  ```
- **Result:** AI compared. Sorted Set gives exact counts and supports per-member expiration directly via score, but memory grows linearly with users per segment. HyperLogLog gives approximate counts with fixed memory (~12KB) regardless of segment size, but needs a day-bucket wrapper for expiration since it has no native remove operation
- **Evaluation:** I read the diffs and according the simple architecture select Sorted Set. But for production-ready architecture HyperLogLog is good fit but is a bit complex to implement.
- **Fixes:** I used Sorted Set for implementation.

### Prompt #3
- **Stage:** Implementation
- **Prompt:**
  ```
  I used background job that periodically fetches expired members in limited batches and removes them
  Is Lua script better way?
  IS my approach cause correctness issues?
  ```
- **Result:** AI suggest Lua script is efficient and confirmed batch approach has no correctness issues. Expired members are already excluded from counts
- **Evaluation:** Trade-off between the two approaches
- **Fixes:** I Keep simple batch, because of keep the implementation simple and readable for the simple architecture.

### Prompt #4
- **Stage:** Refactoring
- **Prompt:**
  ```
  I reviewed the health check implementation
  The handler calling storage ping directly and not using the service layer
  Health check is not just infrastructure concern and we need to ensure application level readiness
  Fix this so the handler only depends on the service and not storage
  ```
- **Result:** AI add `Healthy(ctx context.Context) error` method to `EstimationService` and updated the handler
- **Evaluation:** Reviewing resulting code
- **Fixes:** Change was correct and I accept that.

---

## 4. Token Usage Monitoring

- **Monitoring method:** 
  - built-in context usage, because I used Cursor for developing this project, it shows the context usage in live mode.

- **Management strategies:**
  - I used AI for specific technical questions and did not use it to build the whole project.
  - I did not send the files specially images and did not ask the AI to generate the graphical charts and ASCII diagram text based is good option.
  - I used cheaper or free models for generating simple things like generating documants.
  - I tried to describe anough and give some tips, so in one response get the thing that I need.
