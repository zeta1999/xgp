import time

from gplearn.genetic import SymbolicRegressor
import numpy as np
from sklearn import datasets
from sklearn import ensemble
from sklearn import linear_model
from sklearn import model_selection
from sklearn import pipeline
from sklearn import preprocessing
from sklearn import tree

from xgp import regressor


X, y = datasets.load_boston(return_X_y=True)

cv = model_selection.KFold(n_splits=3, random_state=42)

models = {
    'Random forest': ensemble.RandomForestRegressor(random_state=42),
    'gplearn': SymbolicRegressor(population_size=300,
                           generations=30,
                           p_crossover=0.7, p_subtree_mutation=0.1,
                           p_hoist_mutation=0.05, p_point_mutation=0.1,
                           max_samples=0.9, verbose=0,
                           parsimony_coefficient=0.01, random_state=0),
    'XGP': regressor.XGPRegressor(random_state=42),
    'Lasso': linear_model.Lasso(),
    'Ridge': linear_model.Ridge(),
    'Tree': tree.DecisionTreeRegressor(random_state=42)
}

for name, model in models.items():
    t0 = time.time()
    scores = model_selection.cross_val_score(model, X=X, y=y, scoring='neg_mean_absolute_error', cv=cv)
    print(f'{name}: {-np.mean(scores)} (± {np.std(scores)}) in {time.time() - t0} seconds')
