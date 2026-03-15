import uuid
from fastapi import FastAPI, Request, HTTPException
from contextlib import asynccontextmanager
from pydantic import BaseModel

from sentence_transformers import SentenceTransformer
from qdrant_client import QdrantClient
from qdrant_client.models import VectorParams, Distance, PointStruct
from qdrant_client.http.exceptions import UnexpectedResponse


# 2. Connect to Docker Qdrant
COLLECTION_NAME = "history"

class IndexReq(BaseModel):
    command: str

class IndexResp(BaseModel):
    status: str
    id: str
    

# Initialize the collection if it doesn't exist
@asynccontextmanager
async def lifespan(app: FastAPI):
    # --- STARTUP ---
    print("🚀 Initializing AI Resources...")

    app.state.model = SentenceTransformer('all-MiniLM-L6-v2')
    app.state.q_client = QdrantClient(host="localhost", port=6333)

    if not app.state.q_client.collection_exists(COLLECTION_NAME):
        app.state.q_client.create_collection(
            collection_name=COLLECTION_NAME,
            vectors_config=VectorParams(size=384, distance=Distance.COSINE)
        )

    try:
        yield
    finally:
        print("🛑 Cleaning up resources...")


app = FastAPI(lifespan=lifespan)


@app.post("/index", response_model=IndexResp)
async def index_command(request: Request, body: IndexReq):

    model = request.app.state.model
    q_client = request.app.state.q_client

    try:
        vector = model.encode(body.command).tolist()
        point_id = str(uuid.uuid4())

        q_client.upsert(
            collection_name=COLLECTION_NAME,
            points=[
                PointStruct(
                    id=point_id,
                    vector=vector,
                    payload={"command": body.command}
                )
            ]
        )

        return IndexResp(status="success", id=point_id)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Indexing failed: {str(e)}")



@app.get("/search")
async def search(request: Request, query: str, limit: int = 3):
    
    model = request.app.state.model
    q_client = request.app.state.q_client
    
    try:
        # Convert search query to math
        query_vector = model.encode(query).tolist()
        
        # Find the nearest vectors in the DB
        results = q_client.query_points(
            collection_name=COLLECTION_NAME,
            query=query_vector,
            limit=limit
        ).points
        
        # Return the original commands
        return [{"command": r.payload["command"], "score": r.score} for r in results]
        # return results
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Search failed: {str(e)}")
        